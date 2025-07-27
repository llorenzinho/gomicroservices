package utils

import "golang.org/x/crypto/bcrypt"

func EncryptString(input string) (string, error) {
	hashed, error := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if error != nil {
		return "", error
	}
	return string(hashed), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
