package models

import (
	"log"

	database "github.com/Fabrice-64/golang-architecture/04_client_server/03_jwt/authapp/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

func (user *User) CreateUserRecord() error {
	result := database.GlobalDB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) HashPassword(p string) error {
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		log.Println("error by hashing pwd: ", err)
		return err
	}
	user.Password = string(bs)
	return nil
}

func (user *User) CheckPassword(providedPwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPwd))
	if err != nil {
		log.Println("password does not match: ", err)
		return err
	}
	return nil
}
