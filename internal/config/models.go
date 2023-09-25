package config

// Config represents the configuration options for the application.
type Config struct {
	Notification NotificationConfig `koanf:"notification"`
	App          AppConfig          `koanf:"app"`
}

// NotificationConfig represents the configuration options for notifications.
type NotificationConfig struct {
	Enabled   bool   `koanf:"enabled"`
	GotifyURL string `koanf:"gotify-url"` // Note the - instead of _ to avoid conflicts with env variables
	Username  string `koanf:"username"`
	Password  string `koanf:"password"`
}

// AppConfig represents the configuration options for the application.
type AppConfig struct {
	Port int `koanf:"port"`
}
