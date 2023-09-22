package config

type Config struct {
	Notification NotificationConfig `koanf:"notification"`
	App          AppConfig          `koanf:"app"`
}

type NotificationConfig struct {
	Enabled   bool   `koanf:"enabled"`
	AppToken  string `koanf:"app-token"`
	GotifyURL string `koanf:"gotify-url"` // Note the - instead of _ to avoid conflicts with env variables
}

type AppConfig struct {
	Port int `koanf:"port"`
}
