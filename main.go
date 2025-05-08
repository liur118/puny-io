package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liur/puny-io/internal/config"
	"github.com/liur/puny-io/internal/handler"
	"github.com/liur/puny-io/internal/middleware"
	"github.com/liur/puny-io/internal/service"
	"github.com/liur/puny-io/internal/storage"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化服务
	jwtService := service.NewJWTService(cfg.JwtSecret)
	fileStorage, err := storage.NewFileStorage(cfg.Storage)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// 初始化处理器
	authHandler := handler.NewAuthHandler(jwtService, cfg.Users)
	ossHandler := handler.NewOSSHandler(fileStorage, cfg.Host)

	// 设置路由
	r := gin.Default()

	// 静态文件服务
	// 尝试多个可能的 UI 目录位置
	uiPaths := []string{
		"./ui",                           // 开发环境
		"../ui",                          // 生产环境
		filepath.Join(cfg.Storage, "ui"), // 自定义存储路径
	}

	var uiPath string
	for _, path := range uiPaths {
		if _, err := filepath.Abs(path); err == nil {
			uiPath = path
			break
		}
	}

	if uiPath == "" {
		log.Fatal("Failed to find UI directory")
	}

	// 提供静态文件服务
	r.Static("/ui", uiPath)
	// 将根路径重定向到 UI
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/ui")
	})

	// API 路由组
	api := r.Group("/api")

	// 公开路由
	api.POST("/login", authHandler.Login)
	// 文件访问接口
	r.GET("/oss/:bucket/:key", ossHandler.GetObject)

	// 需要认证的路由
	auth := api.Group("/oss")
	auth.Use(middleware.AuthMiddleware(jwtService))

	// 用户信息
	auth.GET("/user/info", authHandler.GetUserInfo)

	// 桶操作
	auth.GET("/buckets", ossHandler.ListBuckets)
	auth.PUT("/:bucket", ossHandler.CreateBucket)
	auth.DELETE("/:bucket", ossHandler.DeleteBucket)
	auth.GET("/:bucket", ossHandler.ListObjects)

	// 对象操作
	auth.PUT("/:bucket/:key", ossHandler.PutObject)
	auth.DELETE("/:bucket/:key", ossHandler.DeleteObject)
	auth.HEAD("/:bucket/:key", ossHandler.HeadObject)
	auth.PUT("/:bucket/:key/copy", ossHandler.CopyObject)
	auth.GET("/:bucket/:key/url", ossHandler.GetObjectURL)

	// 处理前端路由
	r.NoRoute(func(c *gin.Context) {
		log.Println("NoRoute", c.Request.URL.Path)
		// 检查是否是 API 路径
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}
		// 检查是否是 UI 路径
		if strings.HasPrefix(c.Request.URL.Path, "/ui") {
			// 返回 index.html
			c.File(filepath.Join(uiPath, "index.html"))
			return
		}
		// 其他所有路径都重定向到 /ui
		c.Redirect(302, "/ui")
	})
	// 启动服务器
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
