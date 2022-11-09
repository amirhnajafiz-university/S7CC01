package mqtt

// Config
// rabbitMQ config parameters.
type Config struct {
	Queue string `koanf:"name"`
	URI   string `koanf:"uri"`
}
