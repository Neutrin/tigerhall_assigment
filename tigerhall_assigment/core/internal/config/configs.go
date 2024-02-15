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

	//
	DateTimeFormat = "02-01-2006 15:04:05"
	DateFormat     = "02-01-2006"
)
