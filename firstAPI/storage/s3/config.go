package s3

// Config
// contains data for connecting to s3 database.
type Config struct {
	AccessKeyID     string `koanf:"accessKeyID"`
	SecretAccessKey string `koanf:"secretAccessKey"`
	Region          string `koanf:"region"`
	Bucket          string `koanf:"bucket"`
}
