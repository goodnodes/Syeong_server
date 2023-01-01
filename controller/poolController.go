package controller

import (
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