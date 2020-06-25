package core

// Config defines the configuration for application.
type Config struct {
	// HTTP server address.
	Addr string `json:"addr"`
	// Database
	DB DBConfig `json:"db"`
}

// DBConfig defines configuration for database.
type DBConfig struct {
	// Data source name.
	DSN string `json:"dsn"`
}
