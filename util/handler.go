package util

import (
	"github.com/goodnodes/Syeong_server/log"
)

func ErrorHandler (err error) {
	if err != nil {
		logger.Error(err)
		panic(err)
	}
}