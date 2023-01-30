package model

import (
	// "github.com/goodnodes/Syeong_server/util"
	// "fmt"
	// "encoding/json"
	"context"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagsModel struct {
	Client *mongo.Client
	PoolCollection *mongo.Collection
}

func GetTagsModel(db, host, model string) (*TagsModel, error) {
	tm := &TagsModel{}
	var err error

	if tm.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(host)); err != nil {
		return nil, err
	} else if err = tm.Client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	} else {
		tm.PoolCollection = tm.Client.Database(db).Collection(model)
	}

	return tm, nil
}