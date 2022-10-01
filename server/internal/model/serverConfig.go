package model

type ServerConfig struct {
	SlackConfig    SlackConfig    `yaml:"slack"`
	Linearconfig   LinearConfig   `yaml:"linear"`
	DatabaseConfig DatabaseConfig `yaml:"database"`
	GoogleConfig   GoogleConfig   `yaml:"google"`
	URLConfig      URLConfig      `yaml:"url"`
}

type SlackConfig struct {
	ClientID      string `yaml:"client_id"`
	ClientSecret  string `yaml:"client_secret"`
	SigningSecret string `yaml:"signing_secret"`
	AppToken      string `yaml:"app_token"`
	RedirectURI   string `yaml:"redirect_uri"`
}

type LinearConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURI  string `yaml:"redirect_uri"`
}

type DatabaseConfig struct {
	URI string `yaml:"uri"`
}

type GoogleConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type URLConfig struct {
	App       string `yaml:"app"`
	Dashboard string `yaml:"dashboard"`
}
