package routes

import (
	"golang-chap46/infra"
	"golang-chap46/middleware"

	"github.com/gin-gonic/gin"
)

func NewRoutes(ctx infra.ServiceContext) *gin.Engine {
	router := gin.Default()

	controller := ctx.Ctl
	redisClient := ctx.Cacher.GetRedisClient()
	authMiddleware := ctx.Middleware.Authenticate()

	// Define routes
	router.POST("/register", controller.User.Register)
	router.POST("/login", middleware.RateLimiter(redisClient, 3), controller.User.Login)

	allowedIPs := []string{
		"192.168.1.1", // Example IP
		"203.0.113.0", // Another allowed IP
	}

	productRouter := router.Group("/products", authMiddleware)
	{
		productRouter.GET("/", middleware.IPWhitelistMiddleware(allowedIPs), controller.Product.GetAllProductsController)
	}

	return router
}
