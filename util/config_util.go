package util

import (
	"os"

	"github.com/joho/godotenv"
)

var configurationMap map[string]string

func init() {
	configurationMap = make(map[string]string)

	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		// If .env file doesn't exist, just use environment variables
		Logger.Debugw("No .env file found, using environment variables", "error", err)
	}

	initDefaultValues()
}

func initDefaultValues() {
	// Environment: dev/test/prod
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	configurationMap["env"] = env

	// Logger configuration
	logFolder := os.Getenv("LOG_FOLDER")
	if logFolder == "" {
		logFolder = "./logs"
	}
	configurationMap["logger.file"] = logFolder

	// Database name
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "cashlenx"
	}
	configurationMap["db.name"] = dbName

	// Database type: mongodb / mysql
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "mongodb"
	}
	configurationMap["db.type"] = dbType

	// MongoDB URI format: mongodb+srv://username:password@host/database
	configurationMap["db.mongodb.url"] = os.Getenv("MONGO_DB_URI")

	// MySQL URI format: username:password@tcp(host:port)/database
	configurationMap["db.mysql.url"] = os.Getenv("MYSQL_DB_URI")

	// OpenAPI schema validation: true/false
	schemaValidation := os.Getenv("SCHEMA_VALIDATION")
	if schemaValidation == "" {
		// Enable by default in dev/test environments
		if env == "dev" || env == "test" {
			schemaValidation = "true"
		} else {
			schemaValidation = "false"
		}
	}
	configurationMap["api.schema.validation"] = schemaValidation

	// Authentication configuration
	authEnabled := os.Getenv("AUTH_ENABLED")
	if authEnabled == "" {
		authEnabled = "false" // Disable by default for backward compatibility
	}
	configurationMap["auth.enabled"] = authEnabled

	// JWT Secret for token signing
	jwtSecret := os.Getenv("JWT_SECRET")
	configurationMap["auth.jwt.secret"] = jwtSecret

	// Admin token for sensitive operations
	configurationMap["ADMIN_TOKEN"] = os.Getenv("ADMIN_TOKEN")

	// CORS origins
	corsOrigins := os.Getenv("CORS_ORIGINS")
	configurationMap["cors.origins"] = corsOrigins

	// Log level
	logLevel := os.Getenv("LOG_LEVEL")
	configurationMap["logger.level"] = logLevel

	// Server configuration
	configurationMap["server.port"] = os.Getenv("SERVER_PORT")
	configurationMap["server.host"] = os.Getenv("SERVER_HOST")
	configurationMap["timezone"] = os.Getenv("TIMEZONE")
}

func GetConfigByKey(configKey string) string {
	configValue, isExist := configurationMap[configKey]
	if isExist {
		return configValue
	} else {
		return ""
	}
}

func SetConfigByKey(configKey, configValue string) {
	configurationMap[configKey] = configValue
}
