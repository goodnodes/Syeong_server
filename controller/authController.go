package controller

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	"github.com/goodnodes/Syeong_server/model"
	"github.com/goodnodes/Syeong_server/util"
)

type AuthController struct {
	UserModel *model.UserModel
}

func GetAuthController(um *model.UserModel) *AuthController {
	ac := &AuthController{UserModel : um}

	return ac
}

// 로그인 메서드 - ID(pnum)와 pwd를 받아서 Access/Refresh token을 발급해주는 메서드
func (ac *AuthController) Login(c *gin.Context) {
	var loginStruct model.LoginStruct
	err := c.ShouldBindJSON(&loginStruct)
	util.ErrorHandler(err)

	// req로 받은 pwd를 해시해준다.
	hashedPwd := util.HashPwd(loginStruct.Pwd)

	// user Collection에서 id를 기반으로 user를 찾고, 그 Hpwd와 pwd를 해시한 값이 같은지 비교한다.
	err = util.PwdCompare(loginStruct.Pwd, string(hashedPwd))
	util.ErrorHandler(err)
	// 같다면, RefreshToken과 AccessToken을 발급한다.
	// 이 때 Token의 claims에 들어가는 id는 실제 id가 아닌 db ObjectId로 할 것이다.
	accessToken := util.GetAccessToken(loginStruct.Pnum)
	c.SetCookie("access-token", accessToken, 60*60*24, "/", "localhost:8080", false, true)
	// 여기는 리프레시토큰을 넣어줘야지
	c.SetCookie("refresh-token", accessToken, 60*60*24*60, "/", "localhost:8080", false, true)
	c.JSON(200, gin.H{"msg" : "good"})
}
