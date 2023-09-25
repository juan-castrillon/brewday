package notifications

type Options struct {
	Markdown    bool
	OnClickURL  string
	BigImageURL string
}

type Extras struct {
	Display      *DisplayConfig      `json:"client::display,omitempty"`
	Notification *NotificationConfig `json:"client::notification,omitempty"`
}

type DisplayConfig struct {
	ContentType string `json:"contentType"`
}

type NotificationClick struct {
	Url string `json:"url"`
}

type NotificationConfig struct {
	Click       NotificationClick `json:"click"`
	BigImageURL string            `json:"bigImageUrl"`
}

type Message struct {
	Message  string  `json:"message"`
	Title    string  `json:"title"`
	Priority int     `json:"priority"`
	Extras   *Extras `json:"extras,omitempty"`
}

type ApplicationRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	DefaultPriority int    `json:"defaultPriority,omitempty"`
}

type ApplicationResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
