package util

import (
	"github.com/goodnodes/Syeong_server/log"
)

// error handler util 만들기 -> Logger로 변경 필요
func ErrorHandler (err error) {
	if err != nil {
		logger.Error(err)
		panic(err)
	}
}