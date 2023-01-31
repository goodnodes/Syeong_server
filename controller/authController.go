package controller

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/goodnodes/Syeong_server/model"
	"github.com/goodnodes/Syeong_server/util"
	"io"
	"encoding/json"
	"io/ioutil"
	"github.com/goodnodes/Syeong_server/log"
)

type AuthController struct {
	UserModel *model.UserModel
	ReviewModel *model.ReviewModel
	TagsModel *model.TagsModel
}

func GetAuthController(um *model.UserModel, rm *model.ReviewModel, tm *model.TagsModel) *AuthController {
	ac := &AuthController{
		UserModel : um,
		ReviewModel : rm,
		TagsModel : tm,
	}

	return ac
}

// 로그인 메서드 - ID(pnum)와 pwd를 받아서 Access/Refresh token을 발급해주는 메서드
func (ac *AuthController) Login(c *gin.Context) {
	var loginStruct model.LoginStruct
	err := c.ShouldBindJSON(&loginStruct)
	util.ErrorHandler(err)

	// user Collection에서 id를 기반으로 user를 찾고, 그 Hpwd와 pwd를 해시한 값이 같은지 비교한다.
	user, err := ac.UserModel.FindUserByPnum(loginStruct.Pnum)
	// 존재하지 않는 아이디.
	if err != nil {
		logger.Error(err.Error())
		c.JSON(401, gin.H{
			"msg" : err.Error(), // no information
		})
		return
	}

	// id는 존재하는데, 비밀번호를 전송하기 전
	if loginStruct.Pwd == "" {
		c.JSON(200, gin.H{
			"msg" : "valid id",
		})
		return
	}

	err = util.PwdCompare(user.PrivateInfo.Password, loginStruct.Pwd)
	// 비밀번호가 틀렸을 때
	if err != nil {
		logger.Error(err.Error())
		c.JSON(401, gin.H{
			"err" : "invalid",
		})
		return
	}

	// 같다면, RefreshToken과 AccessToken을 발급한다.
	// 이 때 Token의 claims에 들어가는 id는 실제 id가 아닌 db ObjectId로 할 것이다.
	accessToken := util.GetAccessToken(user.ID.Hex(), user.PrivateInfo.NickName)
	refreshToken := util.GetRefreshToken(user.ID.Hex(), user.PrivateInfo.NickName)
	c.SetCookie("access-token", accessToken, 60*60*24, "/", c.ClientIP(), false, true)
	// c.SetCookie("access-token", accessToken, 10, "/", "localhost", false, true)
	// 여기는 리프레시토큰을 넣어줘야지
	c.SetCookie("refresh-token", refreshToken, 60*60*24*30, "/", c.ClientIP(), false, true)
	// c.SetCookie("refresh-token", refreshToken, 30, "/", "localhost", false, true)
	c.JSON(200, gin.H{"msg" : "good"})
}



// 자동 로그인에서 Access/Refresh token을 확인하는 메서드
// Refresh토큰 유효기간이 얼마 안 남았을 시, Refresh토큰 만료시에 대한 로직
func (ac *AuthController) VerifyToken(c *gin.Context) {
	// RefreshToken 검증작업
	// 여기서 검증이 안되면 에러 반환, 검증 이후 남은 유효기간이 7일 이하면 새 토큰이 발급되어있음
	userId, nickName, err := util.VerifyRefreshToken(c)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(401, gin.H {
			"msg" : err.Error(),
		})
		return
	}

	// 위에서 이미 refreshToken이 검증되었기 때문에, AccessToken을 발급해준다.
	fmt.Println(c.ClientIP())
	fmt.Println(c.Request.UserAgent())
	newAccessToken := util.GetAccessToken(userId, nickName)
	c.SetCookie("access-token", newAccessToken, 60*60*24, "/", c.ClientIP(), false, true)
	// c.SetCookie("access-token", newAccessToken, 10, "/", "localhost", false, true)

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

	purpose := c.Query("purpose")

	pnum := unMarshared["pnum"]

	// 해당 번호로 가입한 사람이 있는지 확인
	result := ac.UserModel.CheckUserByPnum(pnum)
	// 회원가입일 경우, 이미 존재한다면 abort
	if purpose == "signup" {
		if !result {
			c.JSON(401, gin.H{
				"msg" : "already exist",
			})
			return
		}
	} else if purpose == "password" {
		if result {
			c.JSON(401, gin.H{
				"msg" : "not exist",
			})
			return
		}
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
	now := time.Now().Unix()

	err := c.ShouldBindJSON(&check)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	// 인증 시간이 5분을 초과한다면 abort
	if int(now) - check.RequestTime > 300 {
		c.JSON(400, gin.H{
			"msg" : "time over",
		})
		return
	}

	// 문자 requestId를 가지고 문자 content를 확인하여 비교하는 작업
	// 메시지 아이디를 먼저 가져온다.
	messageId := util.GetMsgId(check.RequestId)
	// 메시지 아이디를 사용하여 메시지 콘텐츠를 가져온다.
	messageBody := util.GetMsgContent(messageId)

	// 인증코드만 자른다.
	code := messageBody[17:21]

	// 인증코드 확인 작업
	if check.Code == code {
		verifyCode := util.HashPwd(check.RequestId)
		c.JSON(200, gin.H{
			"msg" : "verified",
			"verifycode" : string(verifyCode),
		})
		return
	} else {
		c.JSON(401, gin.H{
			"msg" : "unverified",
		})
	}
}


// 회원가입 함수 -> 번호 요청하고, 인증 한 후에 진입 가능
func (ac *AuthController) SignUp(c *gin.Context) {
	// verifyCode와 requestId를 가지고 해당 path에 요청할 수 있는지 확인한다.
	verifyCode := c.Query("verifycode")
	requestId := c.Query("requestid")
	err := util.PwdCompare(verifyCode, requestId)
	util.ErrorHandler(err)
	// 먼저 비밀번호와 원하는 닉네임을 받는다.
	user := &model.User{}
	err = c.ShouldBindJSON(user)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	// 다음은 원하는 닉네임이 이미 존재하는지 중복검사를 한다.
	result := ac.UserModel.CheckUserByNickName(user.PrivateInfo.NickName)
	// 이미 있다면 에러 리턴
	if !result {
		c.JSON(401, gin.H{
			"msg" : "already exist",
		})
		return
	}

	// 없다면 회원가입 진행
	// 먼저 createdAt을 채워주자.
	unixTime := time.Now().Unix()
	t := time.Unix(unixTime, 0)
	timeString := t.Format("2006-01-02 15:04:05")
	user.PrivateInfo.CreatedAt = timeString

	// 다음은 패스워드를 해시화해서 저장해주자
	hashedPwd := util.HashPwd(user.PrivateInfo.Password)
	user.PrivateInfo.Password = string(hashedPwd)

	// DB에 유저 정보를 넣어주자
	id, err := ac.UserModel.AddUserData(user)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"id" : id,
	})
}


// 로그아웃 함수 -> 프론트에 전송되었던 쿠키 정보를 삭제
func (ac *AuthController) Logout(c *gin.Context) {
	c.SetCookie("access-token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh-token", "", -1, "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"msg" : "logout success",
	})
}



// 회원탈퇴 함수 -> 회원의 계정을 삭제하고, 연관된 모든 Data를 삭제
// 현재는 리뷰 등 다른 데이터와 연관된 것이 없기 때문에 user만 삭제하면 됨
func (ac *AuthController) DeleteUser(c *gin.Context) {
	userIdString := c.MustGet("userid")
	userId := util.StringToObjectId(userIdString.(string))

	// 먼저 유저가 작성한 리뷰를 모두 가져온다.
	reviews, err := ac.ReviewModel.GetUserReview(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	// 유저가 작성한 리뷰가 있을 경우에만 실행
	if len(reviews) > 0 {
		// 리뷰마다 돌며 태그로 카운트 된 숫자를 감소해준다.
		for _, review := range reviews {
			if len(review.KeywordReviews) == 0 {
				continue
			}
			decTagsArr := util.GetDecTags(review.KeywordReviews...)
			ac.TagsModel.UpdateTagsCount(review.PoolId, decTagsArr)
		}

		// 리뷰를 모두 지워준다.
		err = ac.ReviewModel.DeleteMyReviews(userId)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(400, gin.H{
				"err" : err.Error(),
			})
			return
		}
	}	

	err = ac.UserModel.DeleteMyAccount(userId)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg" : "success",
	})
}



// 비밀번호 변경 함수 -> 번호 요청하고, 인증 한 후에 진입 가능
func (ac *AuthController) ChangePassword(c *gin.Context) {
	// verifyCode와 requestId를 가지고 해당 path에 요청할 수 있는지 확인한다.
	verifyCode := c.Query("verifycode")
	requestId := c.Query("requestid")
	err := util.PwdCompare(verifyCode, requestId)
	util.ErrorHandler(err)
	// 변경하고자 하는 비밀번호를 받는다.
	body := c.Request.Body
	dataMap := make(map[string]interface{})

	data, err := ioutil.ReadAll(body)
	util.ErrorHandler(err)

	json.Unmarshal(data, &dataMap)
	pnum := dataMap["pnum"].(string)
	pwd := dataMap["password"].(string)

	// 비밀번호 변경 진행
	// 패스워드를 해시화 해주자
	hashedPwd := util.HashPwd(pwd)
	
	// pnum을 기준으로 user를 찾아서 해시된 패스워드로 업데이트해주자
	err = ac.UserModel.ChangePassword(pnum, string(hashedPwd))

	if err != nil {
		logger.Error(err.Error())
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg" : "success",
	})
}