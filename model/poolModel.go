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

type PoolModel struct {
	Client *mongo.Client
	PoolCollection *mongo.Collection
}

type Pool struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	City string `bson:"city" json:"city"`
	Region string `bson:"region" json:"region"`
	Name string `bson:"name" json:"name"`
	Url string `bson:"url" json:"url"`
	Address string `bson:"address" json:"address"`
	Pnum string `bson:"pnum" json:"pnum"`
	ImgUrl string `bson:"imgurl" json:"imgurl"`
	OutsideImgUrl string `bson:"outsideimgurl" json:"outsideimgurl"`
	LaneLength int `bson:"lanelength" json:"lanelength"`
	LaneNum int `bson:"lanenum" json:"lanenum"`
	CostInfoUrl string `bson:"costinfourl" json:"costinfourl"`
	FreeSwimInfoUrl string `bson:"freeswiminfourl" json:"freeswiminfourl"`
}

func GetPoolModel(db, host, model string) (*PoolModel, error) {
	pm := &PoolModel{}
	var err error

	if pm.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(host)); err != nil {
		return nil, err
	} else if err = pm.Client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	} else {
		pm.PoolCollection = pm.Client.Database(db).Collection(model)
	}

	return pm, nil
}

func (pm *PoolModel) InsertManyPool(pools []interface{}) error {
	_, err := pm.PoolCollection.InsertMany(context.TODO(), pools)
	
	if err != nil {
		return err
	}
	return nil
}

func (pm *PoolModel) ReplacePool(pool *Pool) error {
	filter := bson.D{{Key : "name", Value : pool.Name}}
	_, err := pm.PoolCollection.ReplaceOne(context.TODO(), filter, pool)
	if err != nil {
		return err
	}
	return nil
}