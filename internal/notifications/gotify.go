package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type GotifyNotifier struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

func NewGotifyNotifier(gotifyURL, username, password string) (*GotifyNotifier, error) {
	n := &GotifyNotifier{
		httpClient: &http.Client{},
		baseURL:    gotifyURL,
	}
	err := n.initializeApp(username, password)
	if err != nil {
		return nil, err
	}
	return n, nil
}

// InitializeApp initializes a gotify app and fills the app token
// It will first check if the app already exists, if not it will create it. The search is done by app name.
func (n *GotifyNotifier) initializeApp(username, password string) error {
	appPath := fmt.Sprintf("%s/application", n.baseURL)
	appName := "brewday"
	log.Info().Msgf("Initializing gotify app %s", appName)
	// See if the app already exists
	req, err := http.NewRequest("GET", appPath, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got status code %d while getting gotify app", resp.StatusCode)
	}
	var appListResp []ApplicationResponse
	err = json.NewDecoder(resp.Body).Decode(&appListResp)
	if err != nil {
		return err
	}
	for _, app := range appListResp {
		if app.Name == appName {
			log.Info().Msgf("Gotify app %s already exists, fetching token", appName)
			n.token = app.Token
			return nil
		}
	}
	// Create the app
	log.Info().Msgf("Gotify app %s does not exist, creating it", appName)
	appReq := ApplicationRequest{
		Name: appName,
	}
	appReqBody, err := json.Marshal(appReq)
	if err != nil {
		return err
	}
	req, err = http.NewRequest("POST", appPath, bytes.NewReader(appReqBody))
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")
	resp, err = n.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got status code %d while creating gotify app", resp.StatusCode)
	}
	var appResp ApplicationResponse
	err = json.NewDecoder(resp.Body).Decode(&appResp)
	if err != nil {
		return err
	}
	n.token = appResp.Token
	log.Info().Msgf("Gotify app %s created", appName)
	return nil
}

// getExtras returns the extras struct for the given options
// This struct is used to configure the notification in the gotify api
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

// SendGotify sends a message to gotify with the given message and title
// It will use the options to configure the notification
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
	messageURL := fmt.Sprintf("%s/message?token=%s", n.baseURL, n.token)
	req, err := http.NewRequest("POST", messageURL, bytes.NewReader(body))
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

// Send sends a message to gotify with the given message and title
// It will use the options to configure the notification
// Supported options are:
// - markdown: bool
// - onClickURL: string
// - bigImageURL: string
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
