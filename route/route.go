package route

import (
	"github.com/goodnodes/Syeong_server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/goodnodes/Syeong_server/controller"
	"github.com/goodnodes/Syeong_server/log"
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

	r.Use(logger.GinLogger())
	r.Use(logger.GinRecovery(true))
	r.Use(middleware.CORS())

	// swagger route 추가 필요

	authGroup := r.Group("/auth")
	{		
		// 번호인증 요청 -> 여기를 통해 회원가입/비밀번호 변경시 메시지를 요청한다.
		authGroup.POST("/request", p.ctl.Auth.RequestNumber)
		// 번호인증 확인 -> 여기를 통해 코드 입력값을 확인한다.
		authGroup.POST("/check", p.ctl.Auth.CheckNumber)
		// 회원가입 -> 앞의 두 단계를 정상진행하면, 아래로 가게 된다.
		authGroup.POST("/signup", p.ctl.Auth.SignUp)
		// 비밀번호 찾기(변경)
		authGroup.POST("/password", p.ctl.Auth.ChangePassword)
		// 로그인
		authGroup.POST("", p.ctl.Auth.Login)
		// 로그아웃
		authGroup.GET("", p.ctl.Auth.Logout)
		// 자동로그인
		authGroup.GET("/auto", p.ctl.Auth.VerifyToken)
	}

	// 작업 경로로 이동할 때는 미들웨어를 사용한다.
	// 이 때 미들웨어에서는 Accesstoken을 확인 후 없으면 Abort 진행 -> 이에 따른 리다이렉팅을 프론트에서 진행해야함
	userGroup := r.Group("/user").Use(middleware.VerifyAccessToken())
	{
		// 나의 정보 가져오기
		userGroup.GET("", p.ctl.User.GetMyInfo)
		// 나의 목표 추가하기
		userGroup.POST("", p.ctl.User.EditMyInfo)
		// 나의 수영장 추가
		userGroup.GET("/pool", p.ctl.User.AddMyPool)
		// 나의 수영장 제거
		userGroup.DELETE("/pool", p.ctl.User.DeleteMyPool)
		// 회원탈퇴 -> 회원 탈퇴 이후 바로 로그아웃 요청할 것
		userGroup.DELETE("", p.ctl.Auth.DeleteUser)
	}

	poolGroup := r.Group("/pool")
	{
		// 전체 수영장 정보 가져오기
		poolGroup.GET("", p.ctl.Pool.GetAll)
		// 수영장별 리뷰 가져오기
		poolGroup.GET("/review", p.ctl.Review.GetPoolReview)
	}

	reviewGroup := r.Group("/review").Use(middleware.VerifyAccessToken())
	{
		// 리뷰 추가하기
		reviewGroup.POST("", p.ctl.Review.AddReview)
		// 유저가 작성한 리뷰 가져오기
		reviewGroup.GET("/user", p.ctl.Review.GetUserReview)
		// 리뷰 수정하기
		reviewGroup.PATCH("", p.ctl.Review.UpdateReview)
		// 리뷰 삭제하기
		reviewGroup.DELETE("", p.ctl.Review.DeleteReview)
	}

	// admin 그룹을 추가하여 pool 정보 업데이트 등을 다루는 메서드를 만들어야함
	adminGroup := r.Group("/admin").Use(middleware.AdminOnly())
	{
		// 수영장 정보 여러개 집어넣기
		adminGroup.POST("", p.ctl.Pool.InsertManyPool)
		// 수영장 정보 업데이트하기
		adminGroup.PATCH("", p.ctl.Pool.ReplacePool)
		// DB에 geo code 업데이트하는 메서드
		adminGroup.GET("", p.ctl.Pool.GetGEO)
		// geo code 하나만 업데이트 하는 메서드
		adminGroup.GET("/one", p.ctl.Pool.GetOneGEO)
	}

	return r
}