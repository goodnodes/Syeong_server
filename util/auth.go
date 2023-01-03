package util

import (
	bcrypt "golang.org/x/crypto/bcrypt"
)

func HashPwd (pwd string) []byte {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	ErrorHandler(err)
	return hashed
}

func PwdCompare (pwd, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err
}