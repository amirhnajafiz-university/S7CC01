package mongodb

type Config struct {
	Enabled  bool   `koanf:"enabled"`
	Database string `koanf:"database"`
	URI      string `koanf:"uri"`
}
