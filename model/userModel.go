package model

import (
	// "fmt"
	// "encoding/json"
	"context"
	
	// "go.mongodb.org/mongo-driver/bson"
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
		CreatedAt string `bson:"createdat" json:"createdat"`
		RefreshToken string `bson:"refreshtoken" json:"refreshtoken"`
		AccessToken string `bson:"accesstoken" json:"accesstoken"`
	}
	MyPools []Pool `bson:"mypools" json:"mypools"`
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