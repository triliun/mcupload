package logger

type Config struct {
	Environment string `json:"environment" yaml:"environment"`
	Level       string `json:"level" yaml:"level"`
	OutputPath  string `json:"output_path" yaml:"output_path"`
	ErrorPath   string `json:"error_path" yaml:"error_path"`
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Environment: "development",
		Level:       "info",
		OutputPath:  "stdout",
		ErrorPath:   "stderr",
	}
}

// NewProductionConfig returns production configuration
func NewProductionConfig() Config {
	return Config{
		Environment: "production",
		Level:       "info",
		OutputPath:  "/var/log/app/app.log",
		ErrorPath:   "/var/log/app/error.log",
	}
}
