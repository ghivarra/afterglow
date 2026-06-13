package extension

import (
	"ghivarra/afterglow/configuration/database"
	"ghivarra/afterglow/environment"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/joho/godotenv"
)

func LoadEnvExtension() {
	// detect .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	environment.LoadServerEnv()
}

func LoadDbExtension() {
	database.ConnectDB()
}

func LoadTimezoneExtension() {
	os.Setenv("TZ", "UTC")
	timezone, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(err)
	}
	time.Local = timezone
}

func LoadVersionExtension() {
	// read version txt file
	versionBytes, errRead := os.ReadFile(environment.SERVER_PATH + "/version.txt")
	if errRead != nil {
		log.Fatal("Version File Not Found!")
	}

	// remove new line
	versionStr := string(versionBytes)
	versionStr = strings.ReplaceAll(versionStr, "\n", "")

	// set env
	environment.SetAppVersion(versionStr)
}
