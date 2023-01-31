package controller

import (
	"fmt"
	"encoding/json"
	"io/ioutil"

	"github.com/goodnodes/Syeong_server/util"
	"github.com/gin-gonic/gin"
	"github.com/goodnodes/Syeong_server/model"
	"github.com/goodnodes/Syeong_server/log"
)

type UserController struct {
	UserModel *model.UserModel
	ReviewModel *model.ReviewModel
	PoolModel *model.PoolModel
}

func GetUserController(um *model.UserModel, rm *model.ReviewModel, pm *model.PoolModel) *UserController {
	uc := &UserController{UserModel : um, ReviewModel : rm, PoolModel : pm}

	return uc
}


// 나의 수영장 추가하는 메서드
func(uc *UserController) AddMyPool(c *gin.Context) {
	userIdString := c.MustGet("userid")
	poolIdString := c.Query("poolid")

	userId := util.StringToObjectId(userIdString.(string))
	poolId := util.StringToObjectId(poolIdString)

	err := uc.UserModel.AddMyPool(userId, poolId)

	if err != nil {
		logger.Error(err)
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg" : "success",
	})
}

// 나의 수영장 제거하는 메서드
func(uc *UserController) DeleteMyPool(c *gin.Context) {
	userIdString := c.MustGet("userid")
	poolIdString := c.Query("poolid")

	userId := util.StringToObjectId(userIdString.(string))
	poolId := util.StringToObjectId(poolIdString)

	err := uc.UserModel.DeleteMyPool(userId, poolId)

	if err != nil {
		logger.Error(err)
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg" : "success",
	})
}


// 나의 정보를 가져오는 메서드
func(uc *UserController) GetMyInfo(c *gin.Context) {
	userIdString := c.MustGet("userid")

	fmt.Println(userIdString)

	userId := util.StringToObjectId(userIdString.(string))

	result := uc.UserModel.GetMyInfo(userId)

	// 패스워드는 보내지 않음
	result.PrivateInfo.Password = ""

	c.JSON(200, gin.H{
		"result" : result,
	})
}


// 나의 목표, 닉네임을 변경하는 메서드
func(uc *UserController) EditMyInfo(c *gin.Context) {
	userIdString := c.MustGet("userid")
	userId := util.StringToObjectId(userIdString.(string))

	body := c.Request.Body
	dataMap := make(map[string]interface{})

	data, err := ioutil.ReadAll(body)
	util.ErrorHandler(err)

	json.Unmarshal(data, &dataMap)
	legacyNickname := dataMap["legacy"].(string)
	newNickName := dataMap["new"].(string)
	goal := dataMap["goal"].(string)

	// 이전 닉넴과 현재 닉넴이 같을땐 바로 goal만 업데이트함
	if legacyNickname == newNickName {
		err = uc.UserModel.EditMyGoal(goal, userId)

		if err != nil {
			logger.Error(err)
			c.JSON(400, gin.H{
				"err" : err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"msg" : "success",
		})
		return
	}

	// 다를 땐 닉네임 비교 후 goal을 업데이트함
	result := uc.UserModel.CheckUserByNickName(newNickName)
	if !result { // 닉네임이 이미 존재하는 경우
		c.JSON(400, gin.H{
			"msg" : "already exist",
		})
		return
	}

	// 닉네임이 존재하지 않는 경우
	// 닉네임 변경하고 goal도 변경해줌
	err = uc.UserModel.EditMyNickName(userId, newNickName)
	if err != nil {
		logger.Error(err)
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}
	err = uc.UserModel.EditMyGoal(goal, userId)
	if err != nil {
		logger.Error(err)
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg" : "success",
	})
}