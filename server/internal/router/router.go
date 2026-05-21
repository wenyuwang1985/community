package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wenyuwang1985/community/internal/handler"
	"github.com/wenyuwang1985/community/internal/middleware"
)

func Setup(mode string, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, communityHandler *handler.CommunityHandler, postHandler *handler.PostHandler, marketHandler *handler.MarketHandler, chatHandler *handler.ChatHandler, wsHandler *handler.WSHandler, jwtSecret string) *gin.Engine {
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

	// WebSocket 连接入口
	v1.GET("/ws", wsHandler.ServeHTTP)

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

		// 广场动态
		authorized.POST("/posts", postHandler.CreatePost)
		authorized.GET("/posts", postHandler.ListPosts)
		authorized.GET("/posts/:id", postHandler.GetPost)
		authorized.DELETE("/posts/:id", postHandler.DeletePost)
		authorized.POST("/posts/:id/like", postHandler.Like)
		authorized.DELETE("/posts/:id/like", postHandler.Unlike)
		authorized.POST("/posts/:id/comments", postHandler.CreateComment)
		authorized.GET("/posts/:id/comments", postHandler.ListComments)

		// 广场集市
		authorized.POST("/items", marketHandler.CreateItem)
		authorized.GET("/items", marketHandler.ListItems)
		authorized.GET("/items/:id", marketHandler.GetItem)
		authorized.PUT("/items/:id", marketHandler.UpdateItem)
		authorized.PUT("/items/:id/sold", marketHandler.MarkSold)
		authorized.PUT("/items/:id/off", marketHandler.MarkOff)

		// 聊天
		authorized.GET("/conversations", chatHandler.ListConversations)
		authorized.POST("/conversations/private", chatHandler.CreatePrivateConversation)
		authorized.GET("/conversations/:id/messages", chatHandler.ListMessages)
	}

	return r
}
