package main

import (
	"time"
	"net/http"
	"github.com/goodnodes/Syeong_server/model"
	"github.com/goodnodes/Syeong_server/controller"
	"github.com/goodnodes/Syeong_server/route"
	"github.com/goodnodes/Syeong_server/config"
)

var cfg = config.GetConfig("config/config.toml")

func main() {
	// config 설정

	port := cfg.Server.Port
	host := cfg.Server.Host
	dbName := cfg.Server.DBname
	userModel := cfg.DB["user"]["model"]
	reviewModel := cfg.DB["review"]["model"]
	poolModel := cfg.DB["pool"]["model"]
	tagsModel := cfg.DB["tags"]["model"]

	// logger 설정 추가 필요

	// 원래는 환경변수를 config.toml 파일에서 받아와야 함

	um, err := model.GetUserModel(dbName, host, userModel)
	if err != nil {
		panic(err)
	}
	rm, err := model.GetReviewModel(dbName, host, reviewModel)
	if err != nil {
		panic(err)
	}
	pm, err := model.GetPoolModel(dbName, host, poolModel)
	if err != nil {
		panic(err)
	}

	tm, err := model.GetTagsModel(dbName, host, tagsModel)
	if err != nil {
		panic(err)
	}

	controller := controller.GetNewController(um, pm, rm, tm)

	router := route.GetRouter(controller)

	mapi := &http.Server {
		Addr : port,
		Handler : router.Idx(),
		ReadTimeout : 5 * time.Second,
		WriteTimeout : 10 * time.Second,
		MaxHeaderBytes : 1 << 20,
	}

	mapi.ListenAndServe()
}