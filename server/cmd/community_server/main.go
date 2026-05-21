package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wenyuwang1985/community/internal/config"
	"github.com/wenyuwang1985/community/internal/database"
	"github.com/wenyuwang1985/community/internal/handler"
	"github.com/wenyuwang1985/community/internal/router"
	"github.com/wenyuwang1985/community/internal/service"
	"github.com/wenyuwang1985/community/internal/ws"
	"github.com/wenyuwang1985/community/pkg/snowflake"
)

func main() {
	// 加载配置
	cfg, err := config.Load("configs/config.dev.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化 Snowflake ID 生成器
	if err := snowflake.Init(cfg.Snowflake.NodeID); err != nil {
		log.Fatalf("初始化 Snowflake 失败: %v", err)
	}
	log.Println("Snowflake ID 生成器初始化成功")

	// 连接 PostgreSQL
	pgPool, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatalf("连接 PostgreSQL 失败: %v", err)
	}
	defer pgPool.Close()
	log.Println("PostgreSQL 连接成功")

	// 连接 Redis
	rdb, err := database.NewRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("连接 Redis 失败: %v", err)
	}
	defer rdb.Close()
	log.Println("Redis 连接成功")

	// 初始化 service 和 handler
	authService := service.NewAuthService(pgPool)
	authHandler := handler.NewAuthHandler(authService, cfg.JWT)

	userService := service.NewUserService(pgPool)
	userHandler := handler.NewUserHandler(userService)

	chatService := service.NewChatService(pgPool)

	communityService := service.NewCommunityService(pgPool)
	communityHandler := handler.NewCommunityHandler(communityService, chatService)

	postService := service.NewPostService(pgPool)
	postHandler := handler.NewPostHandler(postService)

	marketService := service.NewMarketService(pgPool)
	marketHandler := handler.NewMarketHandler(marketService)

	chatHandler := handler.NewChatHandler(chatService)

	wsHub := ws.NewHub()
	wsHandler := handler.NewWSHandler(wsHub, chatService, cfg.JWT.Secret)

	uploadService := service.NewUploadService("")
	uploadHandler := handler.NewUploadHandler(uploadService)

	// 注册路由
	r := router.Setup(cfg.Server.Mode, authHandler, userHandler, communityHandler, postHandler, marketHandler, chatHandler, wsHandler, uploadHandler, cfg.JWT.Secret)

	// 启动 HTTP 服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	go func() {
		log.Printf("服务启动，监听端口 %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭失败: %v", err)
	}
	log.Println("服务已停止")
}
