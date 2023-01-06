package main

import (
	"fmt"
	"time"
	"net/http"
	"github.com/goodnodes/Syeong_server/model"
	"github.com/goodnodes/Syeong_server/controller"
	"github.com/goodnodes/Syeong_server/route"
	"github.com/goodnodes/Syeong_server/config"
)

func main() {
	// config 설정 추가 필요
	config := config.GetConfig("config/config.toml")

	port := config.Server.Port
	host := config.Server.Host
	dbName := config.Server.DBname
	userModel := config.DB["users"]["model"]
	reviewModel := config.DB["reviews"]["model"]
	poolModel := config.DB["pools"]["model"]

	fmt.Println(config)
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

	controller := controller.GetNewController(um, pm, rm)

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