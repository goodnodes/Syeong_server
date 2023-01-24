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

type PoolModel struct {
	Client *mongo.Client
	PoolCollection *mongo.Collection
}

type Pool struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	City string `json:"city"`
	Region string `json:"region"`
	Name string `json:"name"`
	Url string `json:"url"`
	Address string `json:"address"`
	Pnum string `json:"pnum"`
	ImgUrl string `json:"imgurl"`
	OutsideImgUrl string `json:"outsideimgurl"`
	LaneLength int `json:"lanelength"`
	LaneNum int `json:"lanenum"`
	CostInfoUrl string `json:"costinfourl"`
	FreeSwimInfoUrl string `json:"freeswiminfourl"`
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

func (pm *PoolModel) UpsertManyPool() {

}