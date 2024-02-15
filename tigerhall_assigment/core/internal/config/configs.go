package config

const (
	//Database configs
	DBHost     = "localhost"
	DBPort     = "5432"
	DBPassword = "postgres"
	DBUser     = "postgres"
	DBName     = "tigers"

	//Auth key
	//Note will be stored in vault for production
	JwtKey              = "secret_key_for_auth"
	ExpiryTimeInMinutes = 5
)
