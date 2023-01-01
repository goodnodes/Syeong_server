package controller

import (
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

func (*UserController) UserTest(c *gin.Context) {
	c.IndentedJSON(200, gin.H{"msg" : "user router"})
}