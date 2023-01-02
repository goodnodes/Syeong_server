package route

import (
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

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 허용할 header 타입에 대해 열거
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Forwarded-For, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 허용할 method에 대해 열거
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (p *Router) Idx() *gin.Engine {
	r := gin.New()

	// logger 미들웨어 추가 필요
	// recovery 미들웨어 추가 필요
	r.Use(CORS())

	// swagger route 추가 필요

	authGroup := r.Group("/auth")
	{
		// 로그인
		authGroup.POST("", p.ctl.Auth ...)
		// 로그아웃
		authGroup.GET("", p.ctl.Auth ...)
		// 탈퇴
		authGroup.DELETE("", p.ctl.Auth ...)
		// 번호인증 요청
		authGroup.GET("/:pnum", p.ctl.Auth ...)
		// 번호인증 확인
		authGroup.POST("/check", p.ctl.Auth ...)
		// 번호인증 후 회원가입
		authGroup.POST("/signup", p.ctl.Auth ...)
		// 자동로그인
		authGroup.GET("/auto", p.ctl.Auth ...)
	}
	userGroup := r.Group("/user")
	{
		// 나의 정보 가져오기
		userGroup.GET("", p.ctl.User ...)
		// 나의 수영장 추가
		userGroup.POST("/pool", p.ctl.User ...)
		// 나의 수영장 제거
		userGroup.PATCH("/pool", p.ctl.User ...)
	}
	poolGroup := r.Group("/pool")
	{
		// 전체 수영장 정보 가져오기
		poolGroup.GET("", p.ctl.Pool ...)
		// 수영장별 리뷰 가져오기
		poolGroup.GET("/:poolid", p.ctl.Pool ...)
	}
	reviewGroup := r.Group("/review")
	{
		// 리뷰 추가하기
		reviewGroup.POST("", p.ctl.Review ...)
		// 리뷰 수정하기
		reviewGroup.PATCH("", p.ctl.Review ...)
		// 리뷰 삭제하기
		reviewGroup.DELETE("/:reviewid", p.ctl.Review ...)
	}

	return r
}