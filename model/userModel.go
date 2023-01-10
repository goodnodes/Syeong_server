package model

import (
	// "encoding/json"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
	Client *mongo.Client
	UserCollection *mongo.Collection
}

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PrivateInfo struct {
		Pnum string `bson:"pnum" json:"pnum"`
		NickName string `bson:"nickname" json:"nickname"`
		Password string `bson:"password" json:"password"`
		CreatedAt string `bson:"createdat" json:"createdat,omitempty"`
		Goal string `bson:"goal" json:"goal,omitempty"`
	} `bson:"privateinfo" json:"privateinfo"`
	MyPools []Pool `bson:"mypools" json:"mypools,omitempty"`
}

func GetUserModel(db, host, model string) (*UserModel, error) {
	um := &UserModel{}

	var err error

	if um.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(host)); err != nil {
		return nil, err
	} else if err = um.Client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	} else {
		um.UserCollection = um.Client.Database(db).Collection(model)
	}

	return um, nil
}


// 핸드폰 번호 중복검사하는 메서드
func (um *UserModel) FindUserByPnum(pNum string) bool {
	filter := bson.D{{Key : "privateinfo.pnum", Value : pNum}}
	count, _ := um.UserCollection.CountDocuments(context.TODO(), filter)
	fmt.Println(count)
	if count > 0 {
		return false
	} else {
		return true
	}
}


// 닉네임 중복검사하는 메서드
func (um *UserModel) FindUserByNickName(nickName string) bool {

	filter := bson.D{{Key : "privateinfo.nickname", Value : nickName}}
	count, err := um.UserCollection.CountDocuments(context.TODO(), filter)
	fmt.Println(count)
	fmt.Println(err)
	if count > 0 {
		return false
	} else {
		return true
	}
}


// 회원가입 메서드
func (um *UserModel) AddUserData(user *User) (interface{}, error) {
	result, err := um.UserCollection.InsertOne(context.TODO(), user)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}