package controller

import (
	"github.com/goodnodes/Syeong_server/util"
	// "fmt"
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