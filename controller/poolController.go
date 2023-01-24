package controller

import (
	"fmt"

	"github.com/goodnodes/Syeong_server/util"
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

func (pc *PoolController) PoolTest(c *gin.Context) {
	c.IndentedJSON(200, gin.H{"msg" : "pool router"})
}

func (pc *PoolController) UpsertManyPool(c *gin.Context) {
	var pools [] model.Pool
	err := c.ShouldBindJSON(&pools)
	util.ErrorHandler(err)

	for _, value := range pools {
		fmt.Println(value)
	}

	fmt.Println(len(pools))
}