package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Fabrice-64/golang-architecture/04_client_server/03_jwt/authapp/models"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	var actualResult models.User
	user := models.User{
		Name:     "Test User",
		Email:    "jwt@mail.com",
		Password: "secret",
	}
	payload, err := json.Marshal(&user)
	assert.NoError(t, err)
	request, err := http.NewRequest("POST", "/api/public/signup", bytes.NewBuffer(payload))
	assert.NoError(t, err)
}
