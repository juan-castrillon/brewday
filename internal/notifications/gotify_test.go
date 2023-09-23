package notifications

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// setupMockServer sets up a mock http server for testing and a notifier connected to it.
func setupMockServer() (*http.ServeMux, *httptest.Server, *GotifyNotifier) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	n := NewGotifyNotifier("test-token", server.URL)
	return mux, server, n
}

// teardownMock closes the mock server and removes the client.
func teardownMock(server *httptest.Server) {
	server.Close()
}

func TestSend(t *testing.T) {
	require := require.New(t)
	type testCase struct {
		Name    string
		Message string
		Title   string
		Options Options
		Error   bool
	}
	type extraDisplay struct {
		ContentType string `json:"contentType"`
	}
	type extraNotificationClick struct {
		Url string `json:"url"`
	}
	type extraNotification struct {
		Click       extraNotificationClick `json:"click"`
		BigImageURL string                 `json:"bigImageUrl"`
	}
	type extras struct {
		Display      extraDisplay      `json:"client::display,omitempty"`
		Notification extraNotification `json:"client::notification,omitempty"`
	}
	type message struct {
		Message  string `json:"message"`
		Title    string `json:"title"`
		Priority int    `json:"priority"`
		Extras   extras `json:"extras,omitempty"`
	}
	testCases := []testCase{
		{
			Name:    "simple",
			Title:   "test-title",
			Message: "test-message",
			Error:   false,
		},
		{
			Name:    "with-markdown",
			Title:   "test-title",
			Message: "test-message",
			Options: Options{
				Markdown: true,
			},
			Error: false,
		},
		{
			Name:    "with-onclick",
			Title:   "test-title",
			Message: "test-message",
			Options: Options{
				OnClickURL: "https://example.com",
			},
			Error: false,
		},
		{
			Name:    "with-bigimage",
			Title:   "test-title",
			Message: "test-message",
			Options: Options{
				BigImageURL: "https://example.com",
			},
			Error: false,
		},
		{
			Name:    "with-all",
			Title:   "test-title",
			Message: "test-message",
			Options: Options{
				Markdown:    true,
				OnClickURL:  "https://example.com",
				BigImageURL: "https://example.image.com",
			},
			Error: false,
		},
		{
			Name:    "Error in sending",
			Title:   "test-title",
			Message: "test-message",
			Error:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mux, server, n := setupMockServer()
			defer teardownMock(server)
			mux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
				token, ok := r.URL.Query()["token"]
				require.True(ok)
				require.Equal("test-token", token[0])
				if tc.Error {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				require.Equal("POST", r.Method)
				require.Equal("application/json", r.Header.Get("Content-Type"))
				// Read the body
				var msg message
				bytes, err := io.ReadAll(r.Body)
				require.NoError(err)
				defer r.Body.Close()
				require.NoError(json.Unmarshal(bytes, &msg))
				require.Equal(tc.Message, msg.Message)
				require.Equal(tc.Title, msg.Title)
				require.Equal(8, msg.Priority)
				require.Equal(tc.Options.Markdown, msg.Extras.Display.ContentType == "text/markdown")
				require.Equal(tc.Options.OnClickURL, msg.Extras.Notification.Click.Url)
				require.Equal(tc.Options.BigImageURL, msg.Extras.Notification.BigImageURL)
				w.Write([]byte("ok"))
			})
			err := n.SendGotify(tc.Message, tc.Title, tc.Options)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}
