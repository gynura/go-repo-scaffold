package entities_test

import (
	"testing"
	"userservice/internal/core/entities"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/require"
)

func TestUser_EncryptAndSetPassword(t *testing.T) {
	user := &entities.User{}

	err := user.EncryptAndSetPassword("securepassword123")
	require.NoError(t, err)

	hashedPassword := user.Password()
	require.NotEmpty(t, hashedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte("securepassword123"))
	require.NoError(t, err)
}

func TestUser_SetPassword_NoEncryption(t *testing.T) {
	user := &entities.User{}

	err := user.SetPassword("anotherpassword456")
	require.NoError(t, err)

	hashedPassword := user.Password()
	require.NotEmpty(t, hashedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte("anotherpassword456"))
	require.Error(t, err)
}

func TestUser_PasswordSwapping(t *testing.T) {
	user := &entities.User{}

	err := user.EncryptAndSetPassword("testpassword")
	if err != nil {
		return
	}
	require.NoError(t, err)
	initialPassword := user.Password()
	err = user.EncryptAndSetPassword("testpassword2")
	if err != nil {
		return
	}
	require.NoError(t, err)
	require.NotEqual(t, initialPassword, user.Password())
}
