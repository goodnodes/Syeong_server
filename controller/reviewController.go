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
	TagsModel *model.TagsModel
}

func GetReviewController(um *model.UserModel, rm *model.ReviewModel, pm *model.PoolModel, tm *model.TagsModel) *ReviewController {
	rc := &ReviewController{UserModel : um, ReviewModel : rm, PoolModel : pm, TagsModel : tm}

	return rc
}



// 리뷰 작성하는 메서드
func (rc *ReviewController) AddReview(c *gin.Context) {
	review := &model.Review{}
	userId := util.StringToObjectId(c.MustGet("userid").(string))
	nickName := c.MustGet("nickname").(string)
	review.UserId = userId
	err := c.ShouldBindJSON(review)
	util.ErrorHandler(err)

	unixTime := time.Now().Unix()
	t := time.Unix(unixTime, 0)
	timeString := t.Format("2006-01-02 15:04:05")
	review.CreatedAt = timeString

	// Keyword reviews를 가지고 각각의 bson.E 객체를 만들어 배열로 리턴하는 함수
	incTagsArr := util.GetIncTags(review.KeywordReviews...)
	
	rc.TagsModel.UpdateTagsCount(review.PoolId, incTagsArr)

	// 키워드리뷰 정렬
	sort.Strings(review.KeywordReviews)

	// 작성자 닉네임 추가
	review.NickName = nickName

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

	// 동시에 top tag도 계산하여 리턴해줌
	topTags, err := rc.TagsModel.GetTopTags(poolId)
	if err != nil {
		c.JSON(400, gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"reviews" : reviews,
		"topTags" : topTags,
	})
}


// 리뷰 수정하는 메서드
func (rc *ReviewController) UpdateReview(c *gin.Context) {
	userId := c.MustGet("userid").(string)
	reviews := []model.Review{}
	c.ShouldBindJSON(reviews)

	// 요청 전송자와 이전 리뷰 작성자가 다르면 abort
	if reviews[0].UserId != util.StringToObjectId(userId) {
		c.JSON(401, gin.H{
			"err" : "invalid request",
		})
		return
	}

	// 이전 리뷰의 tags count를 Decrease
	decTagsArr := util.GetDecTags(reviews[0].KeywordReviews...)
	err := rc.TagsModel.UpdateTagsCount(reviews[0].PoolId, decTagsArr)
	if err != nil {
		c.JSON(500, gin.H{
			"err" : err.Error(),
		})
		return
	}
	
	// 새 키워드 리뷰 정렬
	sort.Strings(reviews[1].KeywordReviews)

	// 수정일자 추가
	unixTime := time.Now().Unix()
	t := time.Unix(unixTime, 0)
	timeString := t.Format("2006-01-02 15:04:05")
	reviews[1].EditDate = timeString

	// 새 리뷰의 tags count를 Increase
	incTagsArr := util.GetIncTags(reviews[1].KeywordReviews...)
	err = rc.TagsModel.UpdateTagsCount(reviews[1].PoolId, incTagsArr)
	if err != nil {
		c.JSON(500, gin.H{
			"err" : err.Error(),
		})
		return
	}

	err = rc.ReviewModel.UpdateReview(&reviews[1])
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


// 리뷰 삭제하는 메서드
func (rc *ReviewController) DeleteReview(c *gin.Context) {
	reviewId := c.Query("reviewid")
	writerId := c.Query("userid")
	userIdString := c.MustGet("userid").(string)

	// 요청을 보낸 사람과 리뷰 작성자가 같은 사람인지 확인
	if writerId != userIdString {
		c.JSON(401, gin.H{
			"err" : "invalid request",
		})
		return
	}

	// 리뷰 삭제
	err := rc.ReviewModel.DeleteReview(util.StringToObjectId(reviewId))

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


// Have to add Review Tags Collection