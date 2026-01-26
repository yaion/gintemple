package router

import (
	"shop/internal/handler"
	"shop/internal/middleware"
	"shop/internal/websocket"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, userHandler *handler.UserHandler, fileHandler *handler.FileHandler, mw *middleware.Middleware, wsHub *websocket.Hub) {
	// Global middleware
	r.Use(mw.Cors())

	// WebSocket Route
	r.GET("/ws", wsHub.HandleWebSocket)

	// File Upload Routes (Example)
	upload := r.Group("/upload")
	{
		upload.POST("/simple", fileHandler.UploadSimple)
		upload.POST("/init", fileHandler.InitiateMultipart)
		upload.POST("/part", fileHandler.UploadPart)
		upload.POST("/complete", fileHandler.CompleteMultipart)

		// Static file serving for local storage (DEV ONLY)
		r.Static("/uploads", "./uploads")
	}

	// Root API Group
	api := r.Group("/api")

	// 1. SaaS Management (SaaS 管理端)
	// 面向平台管理员：管理租户、计费、系统设置等
	registerSaaSRoutes(api, userHandler, mw)

	// 2. E-commerce Admin (电商后台)
	// 面向商家/租户：管理商品、订单、会员、营销等
	registerAdminRoutes(api, userHandler, mw)

	// 3. E-commerce Mall (电商前台)
	// 面向C端消费者：浏览商品、购物车、下单、个人中心等
	registerMallRoutes(api, userHandler, mw)
}

func registerSaaSRoutes(rg *gin.RouterGroup, h *handler.UserHandler, mw *middleware.Middleware) {
	saas := rg.Group("/saas")
	// saas.Use(mw.Auth("admin")) // Example: Platform Admin Auth
	{
		// 示例：SaaS 平台管理接口
		saas.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "saas module ok"})
		})
	}
}

func registerAdminRoutes(rg *gin.RouterGroup, h *handler.UserHandler, mw *middleware.Middleware) {
	admin := rg.Group("/admin")
	// admin.Use(mw.Auth("merchant")) // Example: Merchant Auth
	{
		// 示例：商家后台接口复用 UserHandler
		admin.GET("/users/:id", h.GetUser)
	}
}

func registerMallRoutes(rg *gin.RouterGroup, h *handler.UserHandler, mw *middleware.Middleware) {
	mall := rg.Group("/mall")
	// mall.Use(mw.Auth("user")) // Example: Customer Auth
	{
		// 示例：前台用户注册
		mall.POST("/register", h.Register)
	}
}
