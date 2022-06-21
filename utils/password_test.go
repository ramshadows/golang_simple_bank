package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword1)

	err = CheckPassword(password, hashPassword1)
	require.NoError(t, err)

	// generate a wrong password
	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashPassword1)
	// now we must receive an error that should be equal to the bcrypt error
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())


	// check that if one pass is hashed twice, the hashes should be different
	hashPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword2)
	require.NotEqual(t, hashPassword1, hashPassword2)

}
