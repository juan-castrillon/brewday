package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Notifier struct {
	httpClient *http.Client
	baseURL    string
}

func NewNotifier(appToken, gotifyURL string) *Notifier {
	baseURL := fmt.Sprintf("%s/message?token=%s", gotifyURL, appToken)
	return &Notifier{
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

func getExtras(o Options) Extras {
	var extras Extras
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

func (n *Notifier) Send(message, title string, opts ...Options) error {
	// var extrasStruct Extras
	// if len(opts) != 0 {
	// 	if len(opts) > 1 {
	// 		return fmt.Errorf("only one extras struct is allowed")
	// 	}
	// 	var options Options
	// 	if len(opts) == 1 {
	// 		options = opts[0]
	// 	}
	// 	extrasStruct = getExtras(options)
	// }
	msg := Message{
		Message:  message,
		Title:    title,
		Priority: 8,
		//Extras:   extrasStruct,
	}
	body, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return err
	}
	// Pretty print the JSON
	fmt.Println(string(body))
	req, err := http.NewRequest("POST", n.baseURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = n.httpClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}
