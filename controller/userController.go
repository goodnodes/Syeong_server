package controller

import (
	// "fmt"
	"encoding/json"
	"io/ioutil"

	"github.com/goodnodes/Syeong_server/util"
	"github.com/gin-gonic/gin"
	"github.com/goodnodes/Syeong_server/model"
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

	userId := util.StringToObjectId(userIdString.(string))

	result := uc.UserModel.GetMyInfo(userId)

	// 패스워드는 보내지 않음
	result.PrivateInfo.Password = ""

	c.JSON(200, gin.H{
		"result" : result,
	})
}


// 나의 목표를 추가하는 메서드
func(uc *UserController) EditMyGoal(c *gin.Context) {
	userIdString := c.MustGet("userid")
	userId := util.StringToObjectId(userIdString.(string))

	body := c.Request.Body
	dataMap := make(map[string]interface{})

	data, err := ioutil.ReadAll(body)
	util.ErrorHandler(err)

	json.Unmarshal(data, &dataMap)
	goal := dataMap["goal"].(string)

	err = uc.UserModel.EditMyGoal(goal, userId)

	if err != nil {
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg" : "success",
	})
}