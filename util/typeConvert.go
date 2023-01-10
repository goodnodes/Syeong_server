package util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StringToObjectId(sId string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(sId)
	ErrorHandler(err)

	return id
}