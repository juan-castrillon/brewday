package ha

import "net/http"

type HANotifer struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

// NewHANotifier will initialize a new HA notifier and try to make a request to test itself
// It returns an error if the request fails
func NewHANotifier(haURL, token string) (*HANotifer, error) {
	n := &HANotifer{
		httpClient: &http.Client{},
		baseURL:    haURL,
		token:      token,
	}
	err := n.healtcheck()
	return n, err
}

func (n *HANotifer) healtcheck() error {
	return nil
}

// Send sends a notification
func (n *HANotifer) Send(message, title string, opts map[string]any) error {
	return nil
}
