package mongodb

// Config
// mongodb config parameters.
type Config struct {
	Enabled  bool   `koanf:"enabled"`
	Database string `koanf:"database"`
	URI      string `koanf:"uri"`
}
