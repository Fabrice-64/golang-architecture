package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	err := initDB()
	assert.NoError(t, err)
}
