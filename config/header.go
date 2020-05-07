package config

// AppConfig app config, mapping 'app.[env].toml' file.
type AppConfig struct {
	Version string // the version of app.
	App     struct {
		Port      int // main service listening port
		PProfPort int // pprof listening port
	}

	// Value custom config.
	Value map[string]interface{}
}
