package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	db "github.com/Fabrice-64/golang-architecture/04_client_server/03_jwt/authapp/database"
)

func TestHashPassword(t *testing.T) {
	u := User{
		Password: "secret",
	}
	err := u.HashPassword(u.Password)
	assert.NoError(t, err)
	os.Setenv("passwordHash", u.Password)
}

func TestCreateUserRecord(t *testing.T) {
	var ur User
	err := db.InitDB()
	if err != nil {
		t.Error(err)
	}
	err = db.GlobalDB.AutoMigrate(&User{})
	assert.NoError(t, err)

	u := User{
		Name:     "Test User",
		Email:    "user@test.com",
		Password: os.Getenv("passwordHash"),
	}

	err = u.CreateUserRecord()
	assert.NoError(t, err)

	db.GlobalDB.Where("email=?", u.Email).Find(&ur)
	db.GlobalDB.Unscoped().Delete(&u)
	assert.Equal(t, "Test User", ur.Name)
	assert.Equal(t, "user@test.com", ur.Email)

}

func TestCheckPassword(t *testing.T) {
	hash := os.Getenv("passwordHash")
	u := User{
		Password: hash,
	}
	err := u.CheckPassword("secret")
	assert.NoError(t, err)
}
