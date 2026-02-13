package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-project/internal/service"
	"github.com/valenrio66/be-project/pkg/token"
	"go.uber.org/zap"

	"github.com/valenrio66/be-project/config"
	_ "github.com/valenrio66/be-project/docs"
	"github.com/valenrio66/be-project/internal/api"
	"github.com/valenrio66/be-project/internal/api/handlers"
	"github.com/valenrio66/be-project/internal/db"
	"github.com/valenrio66/be-project/internal/middleware"
	"github.com/valenrio66/be-project/pkg/database"
	"github.com/valenrio66/be-project/pkg/logger"
)

// @title           Marketing Dashboard API
// @version         1.0
// @description     Backend API for Marketing Dashboard.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:3000
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.InitLogger(cfg.Environment)
	defer logger.Sync()

	logger.Info("Starting Marketing Dashboard Backend...")

	dbPool, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer dbPool.Close()

	log.Println("Database connected successfully")

	tokenMaker := token.NewJWTMaker(cfg.JWTSecret)

	queries := db.New(dbPool)
	userService := service.NewUserService(queries, tokenMaker, cfg)
	campaignService := service.NewCampaignService(queries)
	userHandler := handlers.NewUserHandler(userService)
	campaignHandler := handlers.NewCampaignHandler(campaignService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ZapLogger())
	api.SetupRoutes(r, userHandler, campaignHandler, tokenMaker)

	logger.Info("Testing " + cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		logger.Error("Failed to run server", zap.Error(err))
	}
}
