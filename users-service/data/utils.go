package data

import "golang.org/x/crypto/bcrypt"

func hashPassword(plainText string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 14)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(expected, actual string) error {
	return bcrypt.CompareHashAndPassword([]byte(expected), []byte(actual))
}
