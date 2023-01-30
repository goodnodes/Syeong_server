package util

import (
	"github.com/goodnodes/Syeong_server/model"
)

func GetIncTags(args ...string) []model.Tag {
	var tags []model.Tag

	for _, v := range args {
		tags = append(tags, model.Tag{
			Key : v,
			Value : 1,
		})
	}

	return tags
}