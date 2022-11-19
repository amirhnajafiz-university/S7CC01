package mongodb

// Config
// mongodb config parameters.
type Config struct {
	Enable   bool   `koanf:"enable"`
	Database string `koanf:"database"`
	Local    string `koanf:"local"`
	URI      string `koanf:"uri"`
}
