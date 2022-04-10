package password

import (
	"testing"

	"github.com/betNevS/tinybank/pkg/random"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := random.RandomString(6)

	hashedPassword, err := Hash(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = Check(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := random.RandomString(10)
	err = Check(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
