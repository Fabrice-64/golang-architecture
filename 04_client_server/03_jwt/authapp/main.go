package main

import (
	"log"

	"github.com/Fabrice-64/golang-architecture/04_client_server/03_jwt/authapp/database"
	"github.com/Fabrice-64/golang-architecture/04_client_server/03_jwt/authapp/models"
	"github.com/gin-gonic/gin"
)

func setUpRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func main() {
	err := database.InitDB()
	if err != nil {
		log.Fatalln("could not create DB", err)
	}
	database.GlobalDB.AutoMigrate(&models.User{})
	r := setUpRouter()
	r.Run(":8080")
}
