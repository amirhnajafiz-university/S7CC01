package mongodb

// Config
// mongodb config parameters.
type Config struct {
	Database string `koanf:"database"`
	URI      string `koanf:"uri"`
}
