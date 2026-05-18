package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wenyuwang1985/community/internal/middleware"
)

func Setup(mode string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			middleware.Success(c, "pong")
		})
	}

	return r
}
