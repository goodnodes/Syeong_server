package util

import (
	"go.mongodb.org/mongo-driver/bson"
)

func GetIncTags(args ...string) []bson.E {
	var tags []bson.E

	for _, v := range args {
		tags = append(tags, bson.E{v, 1})
	}

	return tags
}