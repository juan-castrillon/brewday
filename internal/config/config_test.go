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
			Name: "YAML complete only file - gotify",
			Path: "yaml/complete_gotify.yaml",
			Env:  map[string]string{},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
				},
				Store: StoreConfig{
					StoreType: "sql",
					Path:      "./bd.sqlite",
				},
				Process: ProcessParameters{
					LauternRestTimeMin: 15,
					RefractometerWCF:   1.00,
				},
			},
			Error: false,
		},
		{
			Name: "YAML complete only file - ha",
			Path: "yaml/complete_ha.yaml",
			Env:  map[string]string{},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: true,
					Type:    "ha",
					Settings: NotificationSettings{
						HAURL:      "http://localhost:8123",
						HAToken:    "letters1234$_%@*",
						HADeviceID: "mydevice",
					},
				},
				Store: StoreConfig{
					StoreType: "sql",
					Path:      "./bd.sqlite",
				},
				Process: ProcessParameters{
					LauternRestTimeMin: 15,
					RefractometerWCF:   1.00,
				},
			},
			Error: false,
		},
		{
			Name: "YAML complete only file - no defaults - gotify",
			Path: "yaml/complete_no_defaults_gotify.yaml",
			Env:  map[string]string{},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
				},
				Store: StoreConfig{
					StoreType: "sql",
					Path:      "./bd.sqlite",
				},
				Process: ProcessParameters{
					LauternRestTimeMin: 10,
					RefractometerWCF:   1.04,
				},
			},
			Error: false,
		},
		{
			Name: "Only env variables - gotify",
			Path: "",
			Env: map[string]string{
				"BREWDAY_NOTIFICATION_ENABLED":                  "true",
				"BREWDAY_NOTIFICATION_TYPE":                     "gotify",
				"BREWDAY_NOTIFICATION_SETTINGS_GOTIFY-USERNAME": "gotify",
				"BREWDAY_NOTIFICATION_SETTINGS_GOTIFY-PASSWORD": "gotify",
				"BREWDAY_NOTIFICATION_SETTINGS_GOTIFY-URL":      "http://localhost:8080",
				"BREWDAY_APP_PORT":                              "8080",
				"BREWDAY_STORE_TYPE":                            "sql",
				"BREWDAY_STORE_PATH":                            "./bd.sqlite",
				"BREWDAY_PROCESS_LAUTERN-REST-TIME-MIN":         "5",
				"BREWDAY_PROCESS_REFRACTOMETER-WCF":             "1.05",
			},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
				},
				Store: StoreConfig{
					StoreType: "sql",
					Path:      "./bd.sqlite",
				},
				Process: ProcessParameters{
					LauternRestTimeMin: 5,
					RefractometerWCF:   1.05,
				},
			},
			Error: false,
		},
		{
			Name: "Only env variables - ha",
			Path: "",
			Env: map[string]string{
				"BREWDAY_NOTIFICATION_ENABLED":               "true",
				"BREWDAY_NOTIFICATION_TYPE":                  "ha",
				"BREWDAY_NOTIFICATION_SETTINGS_HA-DEVICE-ID": "my_device",
				"BREWDAY_NOTIFICATION_SETTINGS_HA-TOKEN":     "t0ken$#",
				"BREWDAY_NOTIFICATION_SETTINGS_HA-URL":       "http://localhost:8123",
				"BREWDAY_APP_PORT":                           "8080",
				"BREWDAY_STORE_TYPE":                         "sql",
				"BREWDAY_STORE_PATH":                         "./bd.sqlite",
				"BREWDAY_PROCESS_LAUTERN-REST-TIME-MIN":      "5",
				"BREWDAY_PROCESS_REFRACTOMETER-WCF":          "1.05",
			},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: true,
					Type:    "ha",
					Settings: NotificationSettings{
						HAURL:      "http://localhost:8123",
						HAToken:    "t0ken$#",
						HADeviceID: "my_device",
					},
				},
				Store: StoreConfig{
					StoreType: "sql",
					Path:      "./bd.sqlite",
				},
				Process: ProcessParameters{
					LauternRestTimeMin: 5,
					RefractometerWCF:   1.05,
				},
			},
			Error: false,
		},
		{
			Name: "YAML complete file and env variables override - gotify",
			Path: "yaml/complete_gotify.yaml",
			Env: map[string]string{
				"BREWDAY_NOTIFICATION_ENABLED": "false",
				"BREWDAY_STORE_TYPE":           "memory",
				"BREWDAY_STORE_PATH":           "",
			},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: false,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
				Process: ProcessParameters{
					LauternRestTimeMin: 15,
					RefractometerWCF:   1.00,
				},
			},
			Error: false,
		},
		{
			Name: "YAML complete file and env variables override - ha",
			Path: "yaml/complete_ha.yaml",
			Env: map[string]string{
				"BREWDAY_NOTIFICATION_ENABLED": "false",
				"BREWDAY_STORE_TYPE":           "memory",
				"BREWDAY_STORE_PATH":           "",
			},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: false,
					Type:    "ha",
					Settings: NotificationSettings{
						HAURL:      "http://localhost:8123",
						HAToken:    "letters1234$_%@*",
						HADeviceID: "mydevice",
					},
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
				Process: ProcessParameters{
					LauternRestTimeMin: 15,
					RefractometerWCF:   1.00,
				},
			},
			Error: false,
		},
		{
			Name: "Incomplete file and env variables - gotify",
			Path: "yaml/incomplete_gotify.yml",
			Env: map[string]string{
				"BREWDAY_NOTIFICATION_ENABLED":                  "true",
				"BREWDAY_NOTIFICATION_TYPE":                     "gotify",
				"BREWDAY_NOTIFICATION_SETTINGS_GOTIFY-USERNAME": "gotify",
				"BREWDAY_NOTIFICATION_SETTINGS_GOTIFY-PASSWORD": "gotify",
			},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
				Process: ProcessParameters{
					LauternRestTimeMin: 15,
					RefractometerWCF:   1.00,
				},
			},
			Error: false,
		},
		{
			Name: "Incomplete file and env variables - ha",
			Path: "yaml/incomplete_ha.yml",
			Env: map[string]string{
				"BREWDAY_NOTIFICATION_ENABLED":               "true",
				"BREWDAY_NOTIFICATION_TYPE":                  "ha",
				"BREWDAY_NOTIFICATION_SETTINGS_HA-URL":       "http://localhost:8123",
				"BREWDAY_NOTIFICATION_SETTINGS_HA-DEVICE-ID": "mydevice",
			},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: true,
					Type:    "ha",
					Settings: NotificationSettings{
						HAURL:      "http://localhost:8123",
						HAToken:    "letters1234$_%@*",
						HADeviceID: "mydevice",
					},
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
				Process: ProcessParameters{
					LauternRestTimeMin: 15,
					RefractometerWCF:   1.00,
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
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
				Process: ProcessParameters{
					LauternRestTimeMin: 15,
					RefractometerWCF:   1.00,
				},
			},
			Error: false,
		},
		{
			Name:  "Missing User - gotify",
			Path:  "yaml/missing_user_gotify.yaml",
			Error: true,
		},
		{
			Name:  "Missing Password - gotify",
			Path:  "yaml/missing_password_gotify.yaml",
			Error: true,
		},
		{
			Name:  "Missing URL - gotify",
			Path:  "yaml/missing_url_gotify.yaml",
			Error: true,
		},
		{
			Name:  "Missing URL - ha",
			Path:  "yaml/missing_url_ha.yaml",
			Error: true,
		},
		{
			Name:  "Missing Token - ha",
			Path:  "yaml/missing_token_ha.yaml",
			Error: true,
		},
		{
			Name:  "Missing Device ID - ha",
			Path:  "yaml/missing_did_ha.yaml",
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
					Enabled: false,
					Type:    "",
					Settings: NotificationSettings{
						GotifyURL:      "",
						GotifyUsername: "",
						GotifyPassword: "",
					},
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
				Process: ProcessParameters{
					LauternRestTimeMin: 15,
					RefractometerWCF:   1.00,
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
		{
			Name: "Complete config no notification",
			Path: "yaml/complete_no_notification.yaml",
			Env:  map[string]string{},
			Expected: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: false,
				},
				Store: StoreConfig{
					StoreType: "sql",
					Path:      "./bd.sqlite",
				},
				Process: ProcessParameters{
					LauternRestTimeMin: 15,
					RefractometerWCF:   1.00,
				},
			},
			Error: false,
		},
		{
			Name:  "Invalid notification type",
			Path:  "yaml/invalid_notification_type.yaml",
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
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
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
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
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
			Name: "Missing GOTIFY User notification disabled",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: false,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyPassword: "gotify",
					},
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: false,
		},
		{
			Name: "Missing GOTIFY User notification enabled",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyPassword: "gotify",
					},
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: true,
		},
		{
			Name: "Missing GOTIFY Password notification disabled",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: false,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
					},
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: false,
		},
		{
			Name: "Missing GOTIFY Password notification enabled",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
					},
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: true,
		},
		{
			Name: "Missing GOTIFY URL notification disabled",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: false,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
				},
				Store: StoreConfig{
					StoreType: "memory",
				},
			},
			Error: false,
		},
		{
			Name: "Missing GOTIFY URL notification enabled",
			Config: Config{
				App: AppConfig{Port: 8080},
				Notification: NotificationConfig{
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
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
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
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
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
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
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
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
					Enabled: true,
					Type:    "gotify",
					Settings: NotificationSettings{
						GotifyURL:      "http://localhost:8080",
						GotifyUsername: "gotify",
						GotifyPassword: "gotify",
					},
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
