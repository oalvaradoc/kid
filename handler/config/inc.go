package config

// Loader defines the interface of the config loader instance
type Loader interface {
	LoadConfig(filePath string) (*ServiceConfigs, error)
}

// LoaderFactory defines the interface of the config loader factory instance
type LoaderFactory interface {
	CreateConfigLoader() Loader
}
