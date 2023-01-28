package util

import (
	"time"
	"errors"

	bcrypt "golang.org/x/crypto/bcrypt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secret = cfg.Token.Secret


// pwd를 해싱하는 함수
func HashPwd (pwd string) []byte {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	ErrorHandler(err)
	return hashed
}



// pwd를 해싱한 값들을 비교하는 함수
func PwdCompare (hash, pwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err
}



// AccessToken을 발급하는 함수 -> 이 녀석은 RefreshToken이 정상적으로 존재할 때만 재발급한다.
func GetAccessToken(userId, nickName string) string {
	claims := jwt.MapClaims{}
	// claim의 "userid"에 user의 ObjectId를 넣어준다.
	claims["userid"] = userId
	claims["nickname"] = nickName
	// 유효기간 하루
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	// claims["exp"] = time.Now().Add(time.Second * 10).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))

	ErrorHandler(err)

	return signedToken
}



// RefreshToken을 발급하는 함수 -> 이 녀석은 id, pwd를 통해 유저가 DB에 있는 인원인 것을 확인했을 때 정상적으로 발급한다.
func GetRefreshToken(userId, nickName string) string {
	claims := jwt.MapClaims{}
	// claim의 "userid"에 user의 ObjectId를 넣어준다.
	claims["userid"] = userId
	claims["nickname"] = nickName
	// 유효기간 한달
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	// claims["exp"] = time.Now().Add(time.Second * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))

	ErrorHandler(err)

	return signedToken
}



// RefreshToken 검증 및 재발급 하는 함수
func VerifyRefreshToken(c *gin.Context) (string, string, error) {
	
	refreshToken, err := c.Request.Cookie("refresh-token")
	// 쿠키에서 토큰이 잘 가져와졌는지 확인
	if err != nil {
		return "", "", errors.New("get Refresh Cookie failed")
	}
	rtValue := refreshToken.Value
	// 토큰에 value가 잘 들어있는지 확인
	if rtValue == "" {
		return "", "", errors.New("refreshToken is None")
	}

	// 토큰 파싱
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(rtValue, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	// RefreshToken이 검증이 되지 않을 경우에는, 다시 로그인이 필요함
	if err != nil {
		return "", "", errors.New("need to login")
	}

	// 만약 RefreshToken의 유효기간이 일주일 이하로 남았다면 재발급
	if claims["exp"].(float64) < float64(time.Now().Add(time.Hour * 24 * 7).Unix()) {
		newRefreshToken := GetRefreshToken(claims["userid"].(string), claims["nickname"].(string))
		c.SetCookie("refresh-token", newRefreshToken, 60*60*24*30, "/", "localhost:8080", false, true)
	}

	// if claims["exp"].(float64) < float64(time.Now().Add(time.Second * 10).Unix()) {
	// 	newRefreshToken := GetRefreshToken(claims["userid"].(string))
	// 	c.SetCookie("refresh-token", newRefreshToken, 30, "/", "localhost:8080", false, true)
	// }

	return claims["userid"].(string), claims["nickname"].(string), nil
}