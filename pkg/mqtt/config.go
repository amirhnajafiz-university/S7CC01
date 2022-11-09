package mqtt

// Config
// rabbitMQ config parameters.
type Config struct {
	Enabled bool   `koanf:"enabled"`
	Queue   string `koanf:"name"`
	URI     string `koanf:"uri"`
}
