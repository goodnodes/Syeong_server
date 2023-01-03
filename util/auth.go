package util

import (
	"time"

	bcrypt "golang.org/x/crypto/bcrypt"
	jwt "github.com/dgrijalva/jwt-go"
)

// pwd를 해싱하는 함수
func HashPwd (pwd string) []byte {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	ErrorHandler(err)
	return hashed
}

// pwd를 해싱한 값들을 비교하는 함수
func PwdCompare (pwd, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err
}

// AccessToken을 발급하는 함수
func GetAccessToken(userId string) string {
	// 추후 config.toml 파일에서 환경변수 가져올 것
	secret := "secret"

	claims := jwt.MapClaims{}
	// 여기서 userId는 ObjectId일 예정이다.
	claims["userid"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))

	ErrorHandler(err)

	return signedToken
}