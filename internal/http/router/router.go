package router

import (
	_ "github.com/IvanKuchsh-600/subtracker/docs"

	"github.com/IvanKuchsh-600/subtracker/internal/http/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           SubTracker API
// @version         1.0
// @description     REST сервис для агрегации данных об онлайн подписках пользователей
// @termsOfService  http://swagger.io/terms/

// @host           localhost:8080
// @BasePath       /api/v1

func NewRouter(subscriptionHandler *handlers.SubscriptionHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	api := router.Group("/api/v1")

	api.POST("/subscriptions", subscriptionHandler.Create)
	api.GET("/subscriptions", subscriptionHandler.List)
	api.GET("/subscriptions/total-cost", subscriptionHandler.GetTotalCost)
	api.GET("/subscriptions/:id", subscriptionHandler.GetByID)
	api.PUT("/subscriptions/:id", subscriptionHandler.Update)
	api.DELETE("/subscriptions/:id", subscriptionHandler.Delete)

	// Swagger документация
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
