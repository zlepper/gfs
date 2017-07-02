package gfs

import "golang.org/x/crypto/bcrypt"

const (
	passwordCost int = 12
)

// Check the given password
func CheckPassword(plaintext, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Creates a new hashed password from the given input
func CreatePassword(plaintext string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), passwordCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), err
}
