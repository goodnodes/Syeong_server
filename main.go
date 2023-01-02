package main

import (
	"time"
	"net/http"
	"github.com/goodnodes/Syeong_server/model"
	"github.com/goodnodes/Syeong_server/controller"
	"github.com/goodnodes/Syeong_server/route"
	
)

func main() {
	// config 설정 추가 필요
	// logger 설정 추가 필요

	// 원래는 환경변수를 config.toml 파일에서 받아와야 함

	um, err := model.GetUserModel("syeong", "mongodb://localhost:27017", "users")
	if err != nil {
		panic(err)
	}
	rm, err := model.GetReviewModel("syeong", "mongodb://localhost:27017", "reviews")
	if err != nil {
		panic(err)
	}
	pm, err := model.GetPoolModel("syeong", "mongodb://localhost:27017", "pm")
	if err != nil {
		panic(err)
	}

	controller := controller.GetNewController(um, pm, rm)

	router := route.GetRouter(controller)

	mapi := &http.Server {
		Addr : "localhost:8080",
		Handler : router.Idx(),
		ReadTimeout : 5 * time.Second,
		WriteTimeout : 10 * time.Second,
		MaxHeaderBytes : 1 << 20,
	}

	mapi.ListenAndServe()
}