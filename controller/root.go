package controller

import (
	"github.com/goodnodes/Syeong_server/model"
)

type Controller struct {
	User *UserController
	Pool *PoolController
	Review *ReviewController
}

func GetNewController (um *model.UserModel, pm *model.PoolModel, rm *model.ReviewModel) *Controller {
	ctl := &Controller{User : GetUserController(um, rm, pm), Pool : GetPoolController(um, rm, pm), Review : GetReviewController(um, rm, pm)} 
	return ctl
}