package model

import (
	"fmt"
	// "encoding/json"
	"context"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagsModel struct {
	Client *mongo.Client
	TagCollection *mongo.Collection
}

func GetTagsModel(db, host, model string) (*TagsModel, error) {
	tm := &TagsModel{}
	var err error

	if tm.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(host)); err != nil {
		return nil, err
	} else if err = tm.Client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	} else {
		tm.TagCollection = tm.Client.Database(db).Collection(model)
	}

	return tm, nil
}


// 태그 Document를 생성하는 메서드
// func (tm *TagModel) InsertTags(tagId primitive.ObjectID, tags []bson.E) {
// 	tm.
// }


// 태그 업데이트 하는 메서드
func (tm *TagsModel) UpdateTags(tagId primitive.ObjectID, tags []bson.E) error {
	filter := bson.D{{
		Key : "_id", Value : tagId,
	}}
	update := bson.D{{
			Key : "$inc", Value : tags,
		}}
	// Upsert를 사용하여 문서가 있다면 업데이트하고 없다면 추가 후 업데이트
	opts := options.FindOneAndUpdate().SetUpsert(true)
	fmt.Println(update)
	result := tm.TagCollection.FindOneAndUpdate(context.TODO(), filter, update, opts)

	return result.Err()
}

