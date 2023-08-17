package factory

// ConfigOptions defines the config of loader
type ConfigOptions struct {
	Version string
}

// ConfigOption sets an optional config for config loader
type ConfigOption func(*ConfigOptions)

// defines the all supported version of loader
var (
	DefaultVersion = "v2"
)

// NewConfigOptions creates a new ConfigOptions
func NewConfigOptions() ConfigOptions {
	return ConfigOptions{
		Version: DefaultVersion,
	}
}

// WithVersion sets the version number of config loader
func WithVersion(version string) ConfigOption {
	return func(options *ConfigOptions) {
		options.Version = version
	}
}
