package model

import (
	// "github.com/goodnodes/Syeong_server/util"
	// "encoding/json"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"errors"
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
	MyPools []primitive.ObjectID `bson:"mypools,omitempty" json:"mypools,omitempty"`
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
func (um *UserModel) CheckUserByPnum(pNum string) bool {
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
func (um *UserModel) CheckUserByNickName(nickName string) bool {

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



// 핸드폰번호로 user를 찾는 메서드
func (um *UserModel) FindUserByPnum(pNum string) (*User, error) {
	var user User
	filter := bson.D{{
		Key : "privateinfo.pnum", Value : pNum,
	}}

	err := um.UserCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no information")
		}
	}

	return &user, nil
}



// 나의 수영장 추가하는 메서드
func (um *UserModel) AddMyPool(userId, poolId primitive.ObjectID) error {
	filter := bson.D{{Key : "_id", Value : userId}}
	update := bson.D{{Key : "$push", Value : bson.D{{
		Key : "mypools", Value : poolId,
	}}}}
	_, err := um.UserCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}


// 나의 수영장 제거하는 메서드
func (um *UserModel) DeleteMyPool(userId, poolId primitive.ObjectID) error {
	filter := bson.D{{Key : "_id", Value : userId}}
	update := bson.D{{Key : "$pull", Value : bson.D{{
		Key : "mypools", Value : poolId,
	}}}}
	_, err := um.UserCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}



// 나의 정보 가져오는 메서드
func (um *UserModel) GetMyInfo(userId primitive.ObjectID) *User {
	user := &User{}
	filter := bson.D{{
		Key : "_id", Value : userId,
	}}
	um.UserCollection.FindOne(context.TODO(), filter).Decode(user)

	return user
}


// 나의 목표를 수정하는 메서드
func (um *UserModel) EditMyGoal(goal string, userId primitive.ObjectID) error {
	filter := bson.D{{
		Key : "_id", Value : userId,
	}}
	update := bson.D{{
		Key : "$set", Value : bson.D{{
			Key : "privateinfo.goal", Value : goal,
		}},
	}}

	_, err := um.UserCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}


// 닉네임을 변경하는 메서드
func (um *UserModel) EditMyNickName(userId primitive.ObjectID, nickName string) error {
	filter := bson.D{{
		Key : "_id", Value : userId,
	}}
	update := bson.D{{
		Key : "$set", Value : bson.D{{
			Key : "privateinfo.nickname", Value : nickName,
		}},
	}}

	_, err := um.UserCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}


// 회원정보를 삭제하는 메서드
func (um *UserModel) DeleteMyAccount(userId primitive.ObjectID) error {
	filter := bson.D{{
		Key : "_id", Value : userId,
	}}
	_, err := um.UserCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		return err
	}

	return nil
}