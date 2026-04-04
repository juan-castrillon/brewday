package ha

// Message data represent additional options to add to a Message
// Most can be found in https://companion.home-assistant.io/docs/notifications/notifications-basic/#general-options
type MessageData struct {
	ClickAction string `json:"clickAction"`
}

// Message represents a notification message as defined by the notify service in HA
type Message struct {
	Message string       `json:"message"`
	Title   string       `json:"title"`
	Data    *MessageData `json:"data"`
}
