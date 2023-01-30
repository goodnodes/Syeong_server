package model

import (
	"fmt"
	"sort"
	// "encoding/json"
	"context"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagsModel struct {
	Client *mongo.Client
	TagsCollection *mongo.Collection
}

func GetTagsModel(db, host, model string) (*TagsModel, error) {
	tm := &TagsModel{}
	var err error

	if tm.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(host)); err != nil {
		return nil, err
	} else if err = tm.Client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	} else {
		tm.TagsCollection = tm.Client.Database(db).Collection(model)
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
	result := tm.TagsCollection.FindOneAndUpdate(context.TODO(), filter, update, opts)

	return result.Err()
}



// 해당 pool의 tag중 상위 2개만 리턴하는 메서드
func (tm *TagsModel) GetTopTags(tagId primitive.ObjectID) ([]TopTags, error) {
	filter := bson.D{{
		Key : "_id", Value : tagId,
	}}

	// 먼저 document를 가져온다.
	var data map[string]interface{}
	err := tm.TagsCollection.FindOne(context.TODO(), filter).Decode(&data)
	// 가져온 document에서 _id 정보는 지워준다
	delete(data, "_id")

	// 같은 id를 가진 tags document가 없으면 nil 반환
	if err != nil {
		return nil, err
	}

	// document를 내림차순으로 정렬
	var keywords []string
	for key := range data {
		keywords = append(keywords, key)
	}

	sort.SliceStable(keywords, func(x, y int) bool {
		return data[keywords[x]].(int32) > data[keywords[y]].(int32)
	})

	// 정렬한 첫번째 데이터의 value가 0이면 nil 반환
	if data[keywords[0]].(int32) == 0 {
		return nil, err
	// 1. 정렬한 첫 번째 데이터가 0이 아니고, 길이가 1일때
	// 2. 정렬한 두 번째 데이터가 0일 때
	// 첫 번째 값만 반환한다.
	} else if len(keywords) == 1 || data[keywords[1]].(int32) == 0 {
		return []TopTags {{Key : keywords[0], Value : int(data[keywords[0]].(int32))}}, err
	}

	return []TopTags{
		{Key : keywords[0], Value : int(data[keywords[0]].(int32))},
		{Key : keywords[1], Value : int(data[keywords[1]].(int32))},
	}, err
}

