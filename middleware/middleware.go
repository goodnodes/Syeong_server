package middleware

import (

	"github.com/gin-gonic/gin"
	// "github.com/goodnodes/Syeong_server/controller"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/goodnodes/Syeong_server/config"
)

var cfg = config.GetConfig("config/config.toml")
var secret = cfg.Token.Secret
var adminSecret = cfg.Server.Admin

// CORS 설정 미들웨어
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 허용할 header 타입에 대해 열거
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Forwarded-For, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 허용할 method에 대해 열거
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}


// AccessToken 검증 및 재발급 하는 미들웨어
func VerifyAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Request.Cookie("access-token")
		// 쿠키에서 토큰이 잘 가져와졌는지 확인
		if err != nil {
			c.JSON(401, gin.H{
				"msg" : "get Access Cookie failed",
			})
			c.Abort()
			return
		}
		atValue := accessToken.Value
		// 토큰에 value가 잘 들어있는지 확인
		if atValue == "" {
			c.JSON(401, gin.H{
				"msg" : "accessToken is None",
			})
			c.Abort()
			return
		}

		// 토큰 파싱
		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(atValue, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			c.JSON(401, gin.H{
				"msg" : "accessToken is None",
			})
			c.Abort()
		}
		// AccessToken이 검증불가한 경우 RefreshToken을 확인하여 재발급 절차를 거친다.
		// if err != nil {
		// 	var ac *controller.AuthController
		// 	ac.VerifyToken(c)
		// }

		c.Set("userid", claims["userid"])
		c.Set("nickname", claims["nickname"])
		c.Next()
	}
}


func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("password")
		if query != adminSecret {
			c.JSON(401, gin.H{
				"msg" : "not valid admin",
			})
			c.Abort()
		}

		c.Next()
	}
}