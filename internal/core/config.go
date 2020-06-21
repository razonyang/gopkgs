package core

// Config defines the configuration for application.
type Config struct {
	// HTTP server address.
	Addr string `json:"addr"`
	// Plugins location.
	Plugins string `json:"plugins"`
	// Database
	DB DBConfig `json:"db"`
}

// DBConfig defines configuration for database.
type DBConfig struct {
	// Data source name.
	DSN string `json:"dsn"`
	// Database driver.
	Driver string `json:"driver"`
}
