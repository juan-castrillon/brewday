package notifications

// Options are the options for the notification
// These map to the supported options as defined in the gotify API
type Options struct {
	Markdown    bool
	OnClickURL  string
	BigImageURL string
}

// Extras are the extras for the notification
// These are defined by the gotify API
type Extras struct {
	Display      *DisplayConfig      `json:"client::display,omitempty"`
	Notification *NotificationConfig `json:"client::notification,omitempty"`
}

type DisplayConfig struct {
	ContentType string `json:"contentType"`
}

// NotificationClick represents the URL to open when the notification is clicked.
type NotificationClick struct {
	Url string `json:"url"`
}

// NotificationConfig represents the configuration options for a notification.
type NotificationConfig struct {
	Click       NotificationClick `json:"click"`
	BigImageURL string            `json:"bigImageUrl"`
}

// Message represents a notification message.
type Message struct {
	Message  string  `json:"message"`
	Title    string  `json:"title"`
	Priority int     `json:"priority"`
	Extras   *Extras `json:"extras,omitempty"`
}

// ApplicationRequest represents a request to create a new application.
type ApplicationRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	DefaultPriority int    `json:"defaultPriority,omitempty"`
}

// ApplicationResponse represents the response from creating a new application.
type ApplicationResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
