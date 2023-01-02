package controller

import (
	// "github.com/gin-gonic/gin"
	"github.com/goodnodes/Syeong_server/model"
)

type AuthController struct {
	UserModel *model.UserModel
}

func GetAuthController(um *model.UserModel) *AuthController {
	ac := &AuthController{UserModel : um}

	return ac
}