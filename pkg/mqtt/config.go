package mqtt

// Config
// rabbitMQ config parameters.
type Config struct {
	Queue string `koanf:"queue"`
	URI   string `koanf:"uri"`
}
