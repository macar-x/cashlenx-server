package util

import "os"

var configurationMap map[string]string

func init() {
	configurationMap = make(map[string]string)
	initDefaultValues()
}

func initDefaultValues() {
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
