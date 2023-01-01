package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/goodnodes/Syeong_server/model"
)

type PoolController struct {
	UserModel *model.UserModel
	ReviewModel *model.ReviewModel
	PoolModel *model.PoolModel
}

func GetPoolController(um *model.UserModel, rm *model.ReviewModel, pm *model.PoolModel) *PoolController {
	pc := &PoolController{UserModel : um, ReviewModel : rm, PoolModel : pm}

	return pc
}

func (*PoolController) PoolTest(c *gin.Context) {
	c.IndentedJSON(200, gin.H{"msg" : "pool router"})
}