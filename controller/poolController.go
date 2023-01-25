package controller

import (
	// "fmt"

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

func (pc *PoolController) InsertManyPool(c *gin.Context) {
	var pools []interface{}
	err := c.ShouldBindJSON(&pools)
	util.ErrorHandler(err)

	err = pc.PoolModel.InsertManyPool(pools)

	if err != nil {
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg" : "Insert success",
	})
}

func (pc *PoolController) ReplacePool(c *gin.Context) {
	pool := &model.Pool{}
	err := c.ShouldBindJSON(pool)
	util.ErrorHandler(err)

	err = pc.PoolModel.ReplacePool(pool)
	
	if err != nil {
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg" : "Replace success",
	})
}


func (pc *PoolController) GetAll(c *gin.Context) {
	pools, err := pc.PoolModel.GetAllPool()
	if err != nil {
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H {
		"pools" : pools,
	})
}