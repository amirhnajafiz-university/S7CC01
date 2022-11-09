package imagga

// Config
// contains all information for connecting to imagga.
type Config struct {
	ApiKey    string `koanf:"api_key"`
	ApiSecret string `koanf:"api_secret"`
}
