package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/goodnodes/Syeong_server/model"
)

type ReviewController struct {
	UserModel *model.UserModel
	ReviewModel *model.ReviewModel
	PoolModel *model.PoolModel
}

func GetReviewController(um *model.UserModel, rm *model.ReviewModel, pm *model.PoolModel) *ReviewController {
	rc := &ReviewController{UserModel : um, ReviewModel : rm, PoolModel : pm}

	return rc
}



// 리뷰 작성하는 함수
func (*ReviewController) AddReview(c *gin.Context) {
	c.IndentedJSON(200, gin.H{"msg" : "review router"})
}