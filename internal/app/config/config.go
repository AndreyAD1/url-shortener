package config

type StartupConfig struct {
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	ShortURLLength  int    `env:"SHORT_URL_LENGTH" envDefault:"10"`
}
