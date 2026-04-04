package ha

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const MOCK_TOKEN = "6rtqGy3TCVE:-y-a![Yy%_Rv%h6Z*X"

func mockAuth(r *http.Request) bool {
	authContent := r.Header.Get("Authorization")
	expected := "Bearer " + MOCK_TOKEN
	return authContent == expected
}

// setupMockServer sets up a mock http server for testing and a notifier connected to it.
func setupMockServer(token string) (*http.ServeMux, *httptest.Server, *HANotifer, error) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc("/api/{$}", func(w http.ResponseWriter, r *http.Request) {
		ok := mockAuth(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message":"API running"}`))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	n, err := NewHANotifier(server.URL, token, "device1")
	return mux, server, n, err
}

// teardownMock closes the mock server and removes the client.
func teardownMock(server *httptest.Server) {
	server.Close()
}

/*
This works
curl -v -X POST -H "Authorization: Bearer $HA_API_TOKEN" -d@oe.json http://<url>/api/services/notify/mobile_app_juan_phone
with body
{
  "message": "hello from pc",
  "title": "A Title",
  "data": {
    "clickAction": "noAction",
    "color": "red"
  }
}

Body and data params from https://companion.home-assistant.io/docs/notifications/notifications-basic/#opening-a-url

Response ok 200 []
Wrong target 400 "400: Bad Request"
Bad JSON 400 {"message":"Data should be valid JSON."}
Missing message 400 400: Bad Request
Missing title: 200 []
Missing data: 200 []
NE Data input: 200 []
Invalid data input: 200 []
*/

func TestSend(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name    string
		Token   string
		Message string
		Title   string
		Opts    map[string]any
		Error   bool
	}{
		{
			Name:    "Normal case",
			Token:   MOCK_TOKEN,
			Message: "Test 1",
			Title:   "Title 1",
			Opts: map[string]any{
				"clickAction": "noAction",
			},
			Error: false,
		},
		{
			Name:    "False token",
			Token:   "token",
			Message: "Test 2",
			Title:   "Title 2",
			Opts: map[string]any{
				"clickAction": "noAction",
			},
			Error: true,
		},
		{
			Name:    "Empty message",
			Token:   MOCK_TOKEN,
			Message: "",
			Title:   "Title 3",
			Opts: map[string]any{
				"clickAction": "noAction",
			},
			Error: true,
		},
		{
			Name:    "Empty title",
			Token:   MOCK_TOKEN,
			Message: "Test 4",
			Title:   "",
			Opts: map[string]any{
				"clickAction": "noAction",
			},
			Error: false,
		},
		{
			Name:    "Empty opts",
			Token:   MOCK_TOKEN,
			Message: "Test 5",
			Title:   "",
			Opts: map[string]any{
				"clickAction": "",
			},
			Error: false,
		},
		{
			Name:    "Nil Opts",
			Token:   MOCK_TOKEN,
			Message: "Test 6",
			Title:   "",
			Opts:    nil,
			Error:   false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mux, server, n, err := setupMockServer(tc.Token)
			if tc.Token != MOCK_TOKEN {
				require.Error(err)
				return
			}
			require.NoError(err)
			defer teardownMock(server)
			mux.HandleFunc("/api/services/notify/mobile_app_device1", func(w http.ResponseWriter, r *http.Request) {
				ok := mockAuth(r)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				require.Equal("POST", r.Method)
				require.Equal("application/json", r.Header.Get("Content-Type"))
				var msg Message
				bytes, err := io.ReadAll(r.Body)
				require.NoError(err)
				defer r.Body.Close()
				err = json.Unmarshal(bytes, &msg)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				require.Equal(tc.Title, msg.Title)
				require.Equal(tc.Message, msg.Message)
				if tc.Opts == nil {
					require.Equal(&MessageData{}, msg.Data)
				} else {
					require.Equal(tc.Opts["clickAction"].(string), msg.Data.ClickAction)
				}
				if msg.Message == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("[]"))
			})
			err = n.Send(tc.Message, tc.Title, tc.Opts)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}
