package config

type Config struct {
	Notification NotificationConfig `koanf:"notification"`
	App          AppConfig          `koanf:"app"`
}

type NotificationConfig struct {
	Enabled   bool   `koanf:"enabled"`
	GotifyURL string `koanf:"gotify-url"` // Note the - instead of _ to avoid conflicts with env variables
	Username  string `koanf:"username"`
	Password  string `koanf:"password"`
}

type AppConfig struct {
	Port int `koanf:"port"`
}
