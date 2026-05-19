package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wenyuwang1985/community/internal/handler"
	"github.com/wenyuwang1985/community/internal/middleware"
)

func Setup(mode string, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, communityHandler *handler.CommunityHandler, jwtSecret string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			middleware.Success(c, "pong")
		})
	}

	// 公开路由
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	// 社区搜索（公开）
	v1.GET("/communities/search", communityHandler.Search)

	// 用户公开资料
	v1.GET("/users/:id", userHandler.GetUserByID)

	// 需要认证的路由
	authorized := v1.Group("")
	authorized.Use(middleware.AuthRequired(jwtSecret))
	{
		authorized.GET("/user/profile", authHandler.Profile)
		authorized.PUT("/user/profile", userHandler.UpdateProfile)

		authorized.POST("/communities/:id/subscribe", communityHandler.Subscribe)
		authorized.DELETE("/communities/:id/subscribe", communityHandler.Unsubscribe)
		authorized.GET("/user/communities", communityHandler.GetUserSubscriptions)
		authorized.PUT("/user/communities/:id/primary", communityHandler.SetPrimary)
	}

	return r
}
