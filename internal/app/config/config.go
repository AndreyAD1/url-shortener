package config

type StartupConfig struct {
	ServerAddress  string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL        string `env:"BASE_URL" envDefault:"localhost:8080"`
	ShortURLLength int    `env:"SHORT_URL_LENGTH" envDefault:"10"`
}
