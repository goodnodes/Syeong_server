package controller

import (
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