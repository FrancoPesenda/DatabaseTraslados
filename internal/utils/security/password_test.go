package security

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_WhenValidPassword_ShouldReturnHash(t *testing.T) {
	hash, err := HashPassword("secret")

	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "secret", hash)
}

func TestHashPassword_WhenPasswordExceedsBcryptLimit_ShouldReturnError(t *testing.T) {
	hash, err := HashPassword(strings.Repeat("a", 73))

	assert.Error(t, err)
	assert.Empty(t, hash)
}

func TestCheckPassword_WhenPasswordMatches_ShouldReturnNil(t *testing.T) {
	hash, _ := HashPassword("secret")

	err := CheckPassword(hash, "secret")

	assert.Nil(t, err)
}

func TestCheckPassword_WhenPasswordDoesNotMatch_ShouldReturnError(t *testing.T) {
	hash, _ := HashPassword("secret")

	err := CheckPassword(hash, "wrong")

	assert.Error(t, err)
}
