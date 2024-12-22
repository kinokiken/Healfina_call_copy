package main

import (
	"Healfina_call/database"
	_ "Healfina_call/docs"
	_ "Healfina_call/routers"
	"log"

	beego "github.com/beego/beego/v2/server/web"
)

// @title Healfina API
// @version 1.0
// @description API documentation for Healfina project
// @host localhost:8080
// @BasePath /
func main() {
	database.InitMongoDB()
	if database.Client == nil {
		log.Fatal("MongoDB client is not initialized!")
	}
	beego.Run()
}
