package controller

import (
	"github.com/goodnodes/Syeong_server/model"
)

type Controller struct {
	User *UserController
	Pool *PoolController
	Review *ReviewController
	Auth *AuthController
}

func GetNewController (um *model.UserModel, pm *model.PoolModel, rm *model.ReviewModel, tm *model.TagsModel) *Controller {
	ctl := &Controller{
			User : GetUserController(um, rm, pm),
			Pool : GetPoolController(um, rm, pm, tm),
			Review : GetReviewController(um, rm, pm, tm),
			Auth : GetAuthController(um, rm),
		} 
	
	return ctl
}