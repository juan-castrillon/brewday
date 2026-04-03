package ha

import (
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
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
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
	n, err := NewHANotifier(server.URL, token)
	return mux, server, n, err
}

// teardownMock closes the mock server and removes the client.
func teardownMock(server *httptest.Server) {
	server.Close()
}

func TestSend(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name    string
		Token   string
		BaseURL string
		Message string
		Title   string
		Error   bool
	}{
		{Name: ""},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mux, server, n, err := setupMockServer(tc.Token)
			require.NoError(err)
			defer teardownMock(server)
			mux.HandleFunc("/api/aaa", func(w http.ResponseWriter, r *http.Request) {
				ok := mockAuth(r)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			})
			err = n.Send(tc.Message, tc.Title, nil)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}
