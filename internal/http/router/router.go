package router

import (
	"github.com/IvanKuchsh-600/subtracker/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func NewRouter(subscriptionHandler *handlers.SubscriptionHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	api := router.Group("/api/v1")

	subscriptions := api.Group("/subscriptions")

	subscriptions.POST("", subscriptionHandler.Create)
	// subscriptions.GET("", subscriptionHandler.List)
	// subscriptions.GET("/total-cost", subscriptionHandler.GetTotalCost)
	// subscriptions.GET("/:id", subscriptionHandler.GetByID)
	subscriptions.PUT("/:id", subscriptionHandler.Update)
	// subscriptions.DELETE("/:id", subscriptionHandler.Delete)

	// Swagger
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
