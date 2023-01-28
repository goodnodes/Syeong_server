package controller

import (
	// "fmt"
	"time"
	"sort"

	"github.com/goodnodes/Syeong_server/util"
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



// 리뷰 작성하는 메서드
func (rc *ReviewController) AddReview(c *gin.Context) {
	review := &model.Review{}
	userId := util.StringToObjectId(c.MustGet("userid").(string))
	review.UserId = userId
	err := c.ShouldBindJSON(review)
	util.ErrorHandler(err)

	unixTime := time.Now().Unix()
	t := time.Unix(unixTime, 0)
	timeString := t.Format("2006-01-02 15:04:05")
	review.CreatedAt = timeString

	sort.Strings(review.KeywordReviews)

	err = rc.ReviewModel.AddReview(review)

	if err != nil {
		c.JSON(400, gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg" : "success",
	})
}


// 유저가 작성한 리뷰 가져오는 메서드
func (rc *ReviewController) GetUserReview(c *gin.Context) {
	userId := util.StringToObjectId(c.MustGet("userid").(string))
	reviews, err := rc.ReviewModel.GetUserReview(userId)

	if err != nil {
		c.JSON(400, gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"reviews" : reviews,
	})
}


// 수영장별 리뷰 가져오는 메서드
func (rc *ReviewController) GetPoolReview(c *gin.Context) {
	poolId := util.StringToObjectId(c.Query("poolid"))
	reviews, err := rc.ReviewModel.GetPoolReview(poolId)

	if err != nil {
		c.JSON(400, gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"reviews" : reviews,
	})
}


