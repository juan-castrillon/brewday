package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GotifyNotifier struct {
	httpClient *http.Client
	baseURL    string
}

func NewGotifyNotifier(appToken, gotifyURL string) *GotifyNotifier {
	baseURL := fmt.Sprintf("%s/message?token=%s", gotifyURL, appToken)
	return &GotifyNotifier{
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

func getExtras(o Options) *Extras {
	extras := &Extras{
		Display:      &DisplayConfig{},
		Notification: &NotificationConfig{},
	}
	if o.Markdown {
		extras.Display.ContentType = "text/markdown"
	}
	if o.OnClickURL != "" {
		extras.Notification.Click.Url = o.OnClickURL
	}
	if o.BigImageURL != "" {
		extras.Notification.BigImageURL = o.BigImageURL
	}
	return extras
}

func (n *GotifyNotifier) SendGotify(message, title string, opts ...Options) error {
	var extrasStruct *Extras
	if len(opts) != 0 {
		if len(opts) > 1 {
			return fmt.Errorf("only one extras struct is allowed")
		}
		var options Options
		if len(opts) == 1 {
			options = opts[0]
		}
		extrasStruct = getExtras(options)
	}
	msg := Message{
		Message:  message,
		Title:    title,
		Priority: 8,
		Extras:   extrasStruct,
	}
	body, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", n.baseURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
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

func (n *GotifyNotifier) Send(message, title string, opts map[string]any) error {
	options := new(Options)
	if opts != nil {
		options = &Options{}
		if markdown, ok := opts["markdown"].(bool); ok {
			options.Markdown = markdown
		}
		if onClickURL, ok := opts["onClickURL"].(string); ok {
			options.OnClickURL = onClickURL
		}
		if bigImageURL, ok := opts["bigImageURL"].(string); ok {
			options.BigImageURL = bigImageURL
		}
	}
	return n.SendGotify(message, title, *options)
}
