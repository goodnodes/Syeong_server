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
	Geo GEO `bson:"geo" json:"geo"`
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



// admin -> 수영장 정보 여러개 입력하는 메서드
func (pm *PoolModel) InsertManyPool(pools []interface{}) error {
	_, err := pm.PoolCollection.InsertMany(context.TODO(), pools)
	
	if err != nil {
		return err
	}
	return nil
}



// admin -> 수영장 정보 대체하는 메서드
func (pm *PoolModel) ReplacePool(pool *Pool) error {
	filter := bson.D{{Key : "name", Value : pool.Name}}
	_, err := pm.PoolCollection.ReplaceOne(context.TODO(), filter, pool)
	if err != nil {
		return err
	}
	return nil
}

func(pm *PoolModel) GetAllPool() ([]Pool, error) {
	filter := bson.D{{}}
	cursor, err := pm.PoolCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var pools []Pool
	if err = cursor.All(context.TODO(), &pools); err != nil {
		return nil, err
	}

	return pools, nil
}


// admin -> geo code 추가하는 메서드
func(pm *PoolModel) UpdateGEO(poolId primitive.ObjectID, geo *GEO) error {
	filter := bson.D{{
		Key : "_id", Value : poolId,
	}}
	update := bson.D{{
		Key : "$set", Value : bson.D{{
			Key : "geo", Value : geo,
		}},
	}}

	_, err := pm.PoolCollection.UpdateOne(context.TODO(), filter, update)
	return err
}