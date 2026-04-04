package ha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HANotifer struct {
	httpClient *http.Client
	baseURL    string
	token      string
	deviceID   string
}

// NewHANotifier will initialize a new HA notifier and try to make a request to test itself
// It returns an error if the request fails
func NewHANotifier(haURL, token, deviceID string) (*HANotifer, error) {
	n := &HANotifer{
		httpClient: &http.Client{},
		baseURL:    haURL,
		token:      token,
		deviceID:   deviceID,
	}
	err := n.healtcheck()
	return n, err
}

func (n *HANotifer) makeRequest(method, endpoint string, body io.Reader) error {
	url := fmt.Sprintf("%s/%s", n.baseURL, endpoint)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+n.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got status code %d", resp.StatusCode)
	}
	return nil
}

func (n *HANotifer) healtcheck() error {
	return n.makeRequest(http.MethodGet, "api/", nil)
}

// Send sends a notification
func (n *HANotifer) Send(message, title string, opts map[string]any) error {
	data := new(MessageData)
	if opts != nil {
		if clickAction, ok := opts["clickAction"].(string); ok {
			data.ClickAction = clickAction
		}
	}
	m := Message{
		Message: message,
		Title:   title,
		Data:    data,
	}
	body, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s_%s", "api/services/notify/mobile_app", n.deviceID)
	return n.makeRequest(http.MethodPost, url, bytes.NewReader(body))
}
