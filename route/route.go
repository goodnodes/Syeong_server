package route

import (
	"github.com/goodnodes/Syeong_server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/goodnodes/Syeong_server/controller"
)

type Router struct {
	ctl *controller.Controller
}

func GetRouter(ctl *controller.Controller) *Router {
	r := &Router{ctl : ctl}

	return r
}

func (p *Router) Idx() *gin.Engine {
	r := gin.New()

	// logger 미들웨어 추가 필요
	// recovery 미들웨어 추가 필요
	r.Use(middleware.CORS())

	// swagger route 추가 필요

	authGroup := r.Group("/auth")
	{
		// 여기 이제 해야 한다. 문자 메시지를 보내는 방법 추가했음
		// 전화번호를 받으면, 이를 대상으로 랜덤한 숫자 메일을 보낸다. -> 이를 임시저장해야 한다.
		// 임시저장한 번호가 일치하면 이후 로직 가능, 다르면 임시저장 파일 삭제하고 abort
		// 시간 단위로 문자를 보낼때의 시간도 같이 저장해놨다가 시간이 오바되면 거절하는 로직도 추가해야함
		
		// 번호인증 요청
		authGroup.GET("/request", p.ctl.Auth.RequestNumber)
		// // 번호인증 확인
		authGroup.POST("/check", p.ctl.Auth.CheckNumber)
		// // 번호인증 후 회원가입
		// authGroup.POST("/signup", p.ctl.Auth ...)
		// 로그인
		authGroup.POST("", p.ctl.Auth.Login)
		// 로그아웃
		// authGroup.GET("", p.ctl.Auth ...)
		// // 탈퇴
		// authGroup.DELETE("", p.ctl.Auth ...)
		// // 자동로그인
		authGroup.GET("/auto", p.ctl.Auth.VerifyToken)
	}

	// 작업 경로로 이동할 때는 미들웨어를 사용한다.
	// 이 때 미들웨어에서는 Accesstoken을 확인 후 없으면 Abort 진행 -> 이에 따른 리다이렉팅을 프론트에서 진행해야함
	userGroup := r.Group("/user").Use(middleware.VerifyAccessToken())
	// {
	// 	// 나의 정보 가져오기
		userGroup.GET("", p.ctl.User.UserTest)
	// 	// 나의 수영장 추가
	// 	userGroup.POST("/pool", p.ctl.User ...)
	// 	// 나의 수영장 제거
	// 	userGroup.PATCH("/pool", p.ctl.User ...)
	// }

	// poolGroup := r.Group("/pool")
	// {
	// 	// 전체 수영장 정보 가져오기
	// 	poolGroup.GET("", p.ctl.Pool ...)
	// 	// 수영장별 리뷰 가져오기
	// 	poolGroup.GET("/:poolid", p.ctl.Pool ...)
	// }

	// reviewGroup := r.Group("/review")
	// {
	// 	// 리뷰 추가하기
	// 	reviewGroup.POST("", p.ctl.Review ...)
	// 	// 리뷰 수정하기
	// 	reviewGroup.PATCH("", p.ctl.Review ...)
	// 	// 리뷰 삭제하기
	// 	reviewGroup.DELETE("/:reviewid", p.ctl.Review ...)
	// }

	return r
}