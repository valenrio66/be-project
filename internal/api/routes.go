package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/valenrio66/be-project/internal/api/handlers"
	"github.com/valenrio66/be-project/internal/middleware"
	"github.com/valenrio66/be-project/pkg/token"
)

func SetupRoutes(r *gin.Engine, userHandler *handlers.UserHandler, campaignHandler *handlers.CampaignHandler, tokenMaker *token.JWTMaker) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(tokenMaker))
		{
			protected.GET("/me", userHandler.GetMe)
			campaigns := protected.Group("/campaigns")
			{
				campaigns.POST("", campaignHandler.Create)
				campaigns.GET("", campaignHandler.List)
			}
		}
	}
}
