package porgs

// AppBootConfig struct holds configuration required at application boot-up.
type AppBootConfig struct {
	// Host to run the web server on
	Host string

	// Port in the host to run the web server on
	Port int

	// DSN (Data Source Name) for the database connection
	DSN string
}
