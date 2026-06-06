package main

import (
	"ghivarra/afterglow/extension"
	"ghivarra/afterglow/server"
)

// @title Afterglow Auto-Backup REST API
// @version 1.0
// @description Afterglow Auto-Backup Contabo Server REST API using Go + Fiber Framework.
// @contact.name Ghivarra Senandika Rushdie
// @contact.email gsenandika@gmail.com
// @BasePath /
// @securityDefinitions.apiKey ApiKey
// @in header
// @name X-API-KEY
func main() {
	// Run/Load Extensions
	extension.LoadEnvExtension()
	extension.LoadDbExtension()
	extension.LoadTimezoneExtension()
	extension.LoadVersionExtension()

	// run server
	server.CreateServer()
}
