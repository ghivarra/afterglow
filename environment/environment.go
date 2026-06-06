package environment

import (
	"ghivarra/afterglow/src/utility"
	"os"
)

var APP_ENV string
var APP_NAME string
var APP_VERSION string
var APP_KEY string
var SERVER_PORT int
var SERVER_PATH string
var DB_PATH string
var ENCRYPTION_KEY string
var SWAGGER_STATUS string
var SWAGGER_USERNAME string
var SWAGGER_PASSWORD string

func LoadServerEnv() {
	// process
	serverPortStr := os.Getenv("SERVER_PORT")

	// set data
	APP_ENV = os.Getenv("APP_ENV")
	APP_NAME = os.Getenv("APP_NAME")
	APP_KEY = os.Getenv("APP_KEY")
	SERVER_PORT = utility.ParseStrToInt(serverPortStr)
	SERVER_PATH = os.Getenv("SERVER_PATH")
	DB_PATH = os.Getenv("DB_PATH")
	ENCRYPTION_KEY = os.Getenv("ENCRYPTION_KEY")
	SWAGGER_STATUS = os.Getenv("SWAGGER_STATUS")
	SWAGGER_USERNAME = os.Getenv("SWAGGER_USERNAME")
	SWAGGER_PASSWORD = os.Getenv("SWAGGER_PASSWORD")
}

func SetAppVersion(version string) {
	APP_VERSION = version
}
