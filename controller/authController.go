package controller

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/goodnodes/Syeong_server/model"
	"github.com/goodnodes/Syeong_server/util"
	"io"
	"encoding/json"
	"fmt"
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


// 문자인증을 요구하는 메서드
func (ac *AuthController) RequestNumber(c *gin.Context) {
	var unMarshared map[string]string
	body := c.Request.Body
	data, _ := io.ReadAll(body)
	json.Unmarshal(data, &unMarshared)

	pnum := unMarshared["pnum"]
	fmt.Println(pnum)
	// 해당 번호로 가입한 사람이 있는지 확인
	result := ac.UserModel.FindUserByPnum(pnum)
	// 이미 존재한다면 abort
	if !result {
		c.JSON(401, gin.H{
			"msg" : "already exist",
		})
		return
	}

	// 문자 전송하기
	requestId := util.SendMsg(pnum)

	c.JSON(200, gin.H{
		// 메시지 requestId와 입력 시간을 Unix 초 형식으로 전송함
		"requestId" : requestId,
		"requestTime" : time.Now().Unix(),
	})
}


// 문자를 검증하는 메서드
func (ac *AuthController) CheckNumber(c *gin.Context) {
	var check model.CheckSMSStruct
	// now := time.Now().Unix()

	err := c.ShouldBindJSON(&check)
	if err != nil {
		panic(err)
	}

	// 인증 시간이 5분을 초과한다면 abort
	// if int(now) - check.RequestTime > 300 {
	// 	c.JSON(400, gin.H{
	// 		"msg" : "time over",
	// 	})
	// 	return
	// }

	// 문자 requestId를 가지고 문자 content를 확인하여 비교하는 작업
	// 메시지 아이디를 먼저 가져온다.
	messageId := util.GetMsgId(check.RequestId)
	// 메시지 아이디를 사용하여 메시지 콘텐츠를 가져온다.
	messageBody := util.GetMsgContent(messageId)

	// 인증코드만 자른다.
	code := messageBody[17:21]

	// 인증코드 확인 작업
	if check.Code == code {
		c.JSON(200, gin.H{
			"msg" : "verified",
		})
		return
	} else {
		c.JSON(401, gin.H{
			"msg" : "unverified ",
		})
	}
}


// 회원가입 함수 -> 번호 요청하고, 인증 한 후에 진입 가능
func (ac *AuthController) SignUp(c *gin.Context) {

}