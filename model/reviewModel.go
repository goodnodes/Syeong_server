package model

import (
	// "fmt"
	// "encoding/json"
	"context"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewModel struct {
	Client *mongo.Client
	ReviewCollection *mongo.Collection
}

type Review struct {
	// 이렇게 하면 json으로 keyword 하고 text만 받아도 제대로 Object id가 생성되는지 확인해봐야함
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PoolId primitive.ObjectID `bson:"poolid" json:"poolid"`
	UserId primitive.ObjectID `bson:"userid" json:"userid"`
	TextReview string `bson:"textreview" json:"textreview"`
	KeywordReviews  []string `bson:"keywordreviews" json:"keywordreviews"`
	EditDate string `bson:"editdate" json:"editdate"`
	CreatedAt string `bson:"createdat" json:"createdat"`
}

// Review를 다루는 model 객체를 만들어 return 해주는 함수
func GetReviewModel(db, host, model string) (*ReviewModel, error) {
	rm := &ReviewModel{}
	var err error

	if rm.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(host)); err != nil {
		return nil, err
	} else if err = rm.Client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	} else {
		rm.ReviewCollection = rm.Client.Database(db).Collection(model)
	}
	return rm, nil
}


// 리뷰 작성하는 메서드
func (rm *ReviewModel) AddReview(review *Review) error {
	_, err := rm.ReviewCollection.InsertOne(context.TODO(), review)

	return err
}


// 유저가 작성한 리뷰 가져오는 메서드
func (rm *ReviewModel) GetUserReview(userId primitive.ObjectID) ([]Review, error) {
	filter := bson.D{{
		Key : "userid", Value : userId,
	}}
	cursor, err := rm.ReviewCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var reviews []Review
	if err = cursor.All(context.TODO(), &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}


// 수영장 리뷰 가져오는 메서드
func (rm *ReviewModel) GetPoolReview(poolId primitive.ObjectID) ([]Review, error) {
	filter := bson.D{{
		Key : "poolid", Value : poolId,
	}}
	cursor, err := rm.ReviewCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var reviews []Review
	if err = cursor.All(context.TODO(), &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}


// 리뷰 업데이트하는 메서드
func (rm *ReviewModel) UpdateReview(review *Review) error {
	filter := bson.D{{
		Key : "_id", Value : review.ID,
	}}
	update := bson.D{{
		Key : "$set", Value : bson.D{{
			Key : "textreview", Value : review.TextReview,
		}, {
			Key : "keywordreviews", Value : review.KeywordReviews,
		}, {
			Key : "editdate", Value : review.EditDate,
		}},
	}}

	_, err := rm.ReviewCollection.UpdateOne(context.TODO(), filter, update)

	return err
}


// 리뷰 삭제하는 메서드
func (rm *ReviewModel) DeleteReview(reviewId primitive.ObjectID) error {
	filter := bson.D{{
		Key : "_id", Value : reviewId,
	}}
	_, err := rm.ReviewCollection.DeleteOne(context.TODO(), filter)
	return err
}