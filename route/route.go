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

	// logger 미들웨어 추가
	// recovery 미들웨어 추가
	r.Use(CORS())

	// swagger route 추가

	userGroup := r.Group("/user")
	{
		userGroup.GET("", p.ctl.User.UserTest)
	}
	reviewGroup := r.Group("/review")
	{
		reviewGroup.GET("", p.ctl.Review.ReviewTest)
	}
	poolGroup := r.Group("/pool")
	{
		poolGroup.GET("", p.ctl.Pool.PoolTest)
	}

	return r
}