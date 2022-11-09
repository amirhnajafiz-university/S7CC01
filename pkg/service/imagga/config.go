package imagga

// Config
// contains parameters for Imagga service.
type Config struct {
	Enabled bool   `koanf:"enabled"`
	URI     string `koanf:"uri"`
}
