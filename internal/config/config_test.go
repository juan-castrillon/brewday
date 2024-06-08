package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	require := require.New(t)
	basePath := "../../test/config/"
	type testCase struct {
		Name     string
		Path     string
		Env      map[string]string
		Expected Config
		Error    bool
	}
	testCases := []testCase{
		{
			Name:  "invalid path",
			Path:  "invalid",
			Env:   map[string]string{},
			Error: true,
		},
		{
			Name: "YAML complete only file",
			Path: "yaml/complete.yaml",
			Env:  map[string]string{},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "sql",
					Path:      "./bd.sqlite",
				},
			},
			Error: false,
		},
		{
			Name: "Only env variables",
			Path: "",
			Env: map[string]string{
				"BREWDAY_NOTIFICATION_ENABLED":    "true",
				"BREWDAY_NOTIFICATION_USERNAME":   "gotify",
				"BREWDAY_NOTIFICATION_PASSWORD":   "gotify",
				"BREWDAY_NOTIFICATION_GOTIFY-URL": "http://localhost:8080",
				"BREWDAY_APP_PORT":                "8080",
				"BREWDAY_STORE_TYPE":              "sql",
				"BREWDAY_STORE_PATH":              "./bd.sqlite",
			},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "sql",
					Path:      "./bd.sqlite",
				},
			},
			Error: false,
		},
		{
			Name: "YAML complete file and env variables override",
			Path: "yaml/complete.yaml",
			Env: map[string]string{
				"BREWDAY_NOTIFICATION_ENABLED": "false",
				"BREWDAY_STORE_TYPE":           "memory",
				"BREWDAY_STORE_PATH":           "",
			},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   false,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: false,
		},
		{
			Name: "Incomplete file and env variables",
			Path: "yaml/incomplete.yml",
			Env: map[string]string{
				"BREWDAY_NOTIFICATION_ENABLED":  "true",
				"BREWDAY_NOTIFICATION_USERNAME": "gotify",
				"BREWDAY_NOTIFICATION_PASSWORD": "gotify",
			},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: false,
		},
		{
			Name:  "Invalid extension",
			Path:  "yaml/invalid.json",
			Error: true,
		},
		{
			Name: "Title case file",
			Path: "yaml/titleCase.yaml",
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: false,
		},
		{
			Name:  "Missing User",
			Path:  "yaml/missing_user.yaml",
			Error: true,
		},
		{
			Name:  "Missing Password",
			Path:  "yaml/missing_password.yaml",
			Error: true,
		},
		{
			Name:  "Missing URL",
			Path:  "yaml/missing_url.yaml",
			Error: true,
		},
		{
			Name:  "Missing Port",
			Path:  "yaml/missing_port.yaml",
			Error: true,
		},
		{
			Name: "Only port",
			Path: "yaml/only_port.yaml",
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   false,
					GotifyURL: "",
					Username:  "",
					Password:  "",
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: false,
		},
		{
			Name:  "Missing sql path",
			Path:  "yaml/missing_sql_path.yaml",
			Error: true,
		},
		{
			Name:  "Missing sql path",
			Path:  "yaml/invalid_store.yaml",
			Error: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			for k, v := range tc.Env {
				err := os.Setenv(k, v)
				require.NoError(err)
				defer os.Unsetenv(k)
			}
			p := ""
			if tc.Path != "" {
				p = filepath.Join(basePath, tc.Path)
			}
			actual, err := LoadConfig(p)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tc.Expected, *actual)
			}
		})
	}
}

func TestFormatEnvVariables(t *testing.T) {
	type testCase struct {
		Name        string
		InputKey    string
		InputVal    string
		ExpectedKey string
	}
	testCases := []testCase{
		{Name: "simple", InputKey: "BREWDAY_NOTIFICATION_ENABLED", InputVal: "true", ExpectedKey: "notification.enabled"},
		{Name: "underscore", InputKey: "BREWDAY_NOTIFICATION_APP_TOKEN", InputVal: "token", ExpectedKey: "notification.app.token"},
		{Name: "hyphen", InputKey: "BREWDAY_NOTIFICATION_GOTIFY-URL", InputVal: "http://localhost:8080", ExpectedKey: "notification.gotify-url"},
		{Name: "noPrefix", InputKey: "NOTIFICATION_GOTIFY-URL", InputVal: "http://localhost:8080", ExpectedKey: "notification.gotify-url"},
		{Name: "otherPrefix", InputKey: "OTHER_NOTIFICATION_GOTIFYURL", InputVal: "http://localhost:8080", ExpectedKey: "other.notification.gotifyurl"},
		{Name: "empty", InputKey: "", InputVal: "http://localhost:8080", ExpectedKey: ""},
		{Name: "noValue", InputKey: "BREWDAY_KEY", InputVal: "", ExpectedKey: "key"},
		{Name: "severalLevels", InputKey: "BREWDAY_LEVEL1_LEVEL2_LEVEL3", InputVal: "value", ExpectedKey: "level1.level2.level3"},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actualKey, actualVal := formatEnvVariables(tc.InputKey, tc.InputVal)
			require.Equal(t, tc.ExpectedKey, actualKey)
			actualValStr, ok := actualVal.(string)
			require.True(t, ok)
			require.Equal(t, tc.InputVal, actualValStr)
		})
	}
}

func TestGetParser(t *testing.T) {
	type testCase struct {
		Name         string
		Format       string
		ExpectedType string
		Error        bool
	}
	testCases := []testCase{
		{Name: "yaml", Format: ".yaml", ExpectedType: "*yaml.YAML", Error: false},
		{Name: "yml", Format: ".yml", ExpectedType: "*yaml.YAML", Error: false},
		{Name: "unsupported", Format: ".json", ExpectedType: "", Error: true},
		{Name: "empty", Format: "", ExpectedType: "", Error: true},
		{Name: "invalid", Format: "invalid", ExpectedType: "", Error: true},
		{Name: "NoDot", Format: "yaml", ExpectedType: "", Error: true},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual, err := getParser(tc.Format)
			if tc.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				actualType := reflect.TypeOf(actual).String()
				require.Equal(t, tc.ExpectedType, actualType)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	type testCase struct {
		Name   string
		Config Config
		Error  bool
	}
	testCases := []testCase{
		{
			Name: "Valid",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: false,
		},
		{
			Name: "Valid SQL",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "sql",
					Path:      "./bd.sqlite",
				},
			},
			Error: false,
		},
		{
			Name: "Missing Port",
			Config: Config{
				App: AppConfig{Port: 0},
			},
			Error: true,
		},
		{
			Name: "Missing User notification disabled",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   false,
					GotifyURL: "http://localhost:8080",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: false,
		},
		{
			Name: "Missing User notification enabled",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "http://localhost:8080",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: true,
		},
		{
			Name: "Missing Password notification disabled",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   false,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: false,
		},
		{
			Name: "Missing Password notification enabled",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: true,
		},
		{
			Name: "Missing URL notification disabled",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   false,
					GotifyURL: "",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: false,
		},
		{
			Name: "Missing URL notification enabled",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: true,
		},
		{
			Name: "Invalid store",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "invalid",
					Path:      "./bd.sqlite",
				},
			},
			Error: true,
		},
		{
			Name: "SQL missing path",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "sql",
				},
			},
			Error: true,
		},
		{
			Name: "SQL Store in uppercase",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "SQL",
					Path:      "./bd.sqlite",
				},
			},
			Error: true,
		},
		{
			Name: "Memory store in uppercase",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled:   true,
					GotifyURL: "http://localhost:8080",
					Username:  "gotify",
					Password:  "gotify",
				},
				Store: StoreConfig{
					StoreType: "MEMORY",
				},
			},
			Error: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := validateConfig(&tc.Config)
			if tc.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
