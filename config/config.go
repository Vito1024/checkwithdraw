package config

type Config struct {
	OkLinkConfig
}

func New() *Config {
	return &Config{
		OkLinkConfig: OkLinkConfig{
			Host:   "https://www.oklink.com",
			ApiKey: "77c8661f-4db1-4cef-81c8-b5e6b5a6dc22",
		},
	}
}

type OkLinkConfig struct {
	Host   string
	ApiKey string
}
