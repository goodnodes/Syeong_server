package main

import (
	"fmt"
	"time"
	"net/http"
	"github.com/goodnodes/Syeong_server/model"
	"github.com/goodnodes/Syeong_server/controller"
	"github.com/goodnodes/Syeong_server/route"
	"github.com/goodnodes/Syeong_server/config"
	"github.com/goodnodes/Syeong_server/log"
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
	if err := logger.InitLogger(cfg); err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}

	logger.Debug("ready server....")

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