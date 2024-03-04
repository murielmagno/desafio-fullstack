package main

import (
	"desafio-fullstack/database"
	"desafio-fullstack/routes"
)

// @title Account API
// @version 1.0
// @description Create Go REST API in Gin Framework
// @termsOfService demo.com
// @contact.url http://demo.com/support
// @contact.email support@swagger.io
// @host localhost:8080
// @BasePath /api
// @Schemes http https
// @query.collection.format multi
// @in header
func main() {
	database.ConnectDB()
	routes.AccountHandler()
}
