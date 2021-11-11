package controllers

import (
	"log"

	"github.com/Fabrice-64/golang-architecture/04_client_server/03_jwt/authapp/models"
	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"msg": "invalid json",
		})
		c.Abort()
		return
	}
	err = user.HashPassword(user.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{
			"msg": "error hashing password",
		})
		c.Abort()
		return
	}
	err = user.CreateUserRecord()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"msg": "error creating useer",
		})
		c.Abort()
		return
	}
	c.JSON(200, user)
}
