package controller

import (
	// "fmt"

	"github.com/goodnodes/Syeong_server/util"
	"github.com/gin-gonic/gin"
	"github.com/goodnodes/Syeong_server/model"
	"github.com/goodnodes/Syeong_server/log"
)

type PoolController struct {
	UserModel *model.UserModel
	ReviewModel *model.ReviewModel
	PoolModel *model.PoolModel
	TagsModel *model.TagsModel
}

func GetPoolController(um *model.UserModel, rm *model.ReviewModel, pm *model.PoolModel, tm *model.TagsModel) *PoolController {
	pc := &PoolController{UserModel : um, ReviewModel : rm, PoolModel : pm, TagsModel : tm}

	return pc
}


// admin -> 여러개의 풀 정보를 한번에 넣는 메서드
func (pc *PoolController) InsertManyPool(c *gin.Context) {
	var pools []interface{}
	err := c.ShouldBindJSON(&pools)
	util.ErrorHandler(err)

	err = pc.PoolModel.InsertManyPool(pools)

	if err != nil {
		logger.Error(err.Error())
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg" : "Insert success",
	})
}


// 풀 정보를 업데이트하는 메서드 (대체)
func (pc *PoolController) ReplacePool(c *gin.Context) {
	pool := &model.Pool{}
	err := c.ShouldBindJSON(pool)
	util.ErrorHandler(err)

	err = pc.PoolModel.ReplacePool(pool)
	
	if err != nil {
		logger.Error(err.Error())
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg" : "Replace success",
	})
}


// 모든 풀 정보를 가져오는 메서드
func (pc *PoolController) GetAll(c *gin.Context) {
	pools, err := pc.PoolModel.GetAllPool()
	if err != nil {
		logger.Error(err.Error())
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H {
		"pools" : pools,
	})
}


// admin -> geo code를 가져와서 업데이트해주는 메서드
func (pc *PoolController) GetGEO(c *gin.Context) {
	// 먼저 전체 수영장 정보를 가져온다
	pools, err := pc.PoolModel.GetAllPool()
	if err != nil {
		logger.Error(err.Error())
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	// fmt.Println(pools)

	// 모든 수영장 요소에 대해서 과정을 진행한다.
	for _, value := range pools {
		geo := util.GetGEO(value.Address)
		err = pc.PoolModel.UpdateGEO(value.ID, geo)
		if err != nil {
			logger.Error(err.Error())
			break
		}
	}

	if err != nil {
		logger.Error(err.Error())
		c.JSON(400, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg" : "geo success",
	})
}