package controller

import (
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
	refreshToken := util.GetRefreshToken(loginStruct.Pnum)
	// c.SetCookie("access-token", accessToken, 60*60*24, "/", "localhost", false, true)
	c.SetCookie("access-token", accessToken, 10, "/", "localhost", false, true)
	// 여기는 리프레시토큰을 넣어줘야지
	// c.SetCookie("refresh-token", refreshToken, 60*60*24*30, "/", "localhost", false, true)
	c.SetCookie("refresh-token", refreshToken, 30, "/", "localhost", false, true)
	c.JSON(200, gin.H{"msg" : "good"})
}



// 자동 로그인에서 Access/Refresh token을 확인하는 메서드
// Refresh토큰 유효기간이 얼마 안 남았을 시, Refresh토큰 만료시에 대한 로직
func (ac *AuthController) VerifyToken(c *gin.Context) {
	// RefreshToken 검증작업
	// 여기서 검증이 안되면 에러 반환, 검증 이후 남은 유효기간이 7일 이하면 새 토큰이 발급되어있음
	userId, err := util.VerifyRefreshToken(c)
	if err != nil {
		c.JSON(401, gin.H {
			"msg" : err.Error(),
		})
		return
	}

	// 위에서 이미 refreshToken이 검증되었기 때문에, AccessToken을 발급해준다.
	newAccessToken := util.GetAccessToken(userId)
	// c.SetCookie("access-token", newAccessToken, 60*60*24, "/", "localhost", false, true)
	c.SetCookie("access-token", newAccessToken, 10, "/", "localhost", false, true)

	c.JSON(200, gin.H {
		"msg" : "success",
	})
}

func (ac *AuthController) SignUp(c *gin.Context) {

}