package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gptscript-ai/knowledge/pkg/db"
	"github.com/gptscript-ai/knowledge/pkg/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Config struct {
	ServerURL, Port, APIBase string
}

type Server struct {
	db *db.DB
}

func NewServer(db *db.DB) *Server {
	return &Server{db: db}
}

// Start starts the server with the given configuration.
func (s *Server) Start(ctx context.Context, cfg Config) error {
	router := gin.Default()

	// Database migration
	if err := s.db.AutoMigrate(); err != nil {
		return err
	}

	// API routes
	docs.SwaggerInfo.BasePath = cfg.APIBase

	// @title Knowledge API
	v1 := router.Group(cfg.APIBase)
	{
		v1.POST("/datasets/create", CreateDataset)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // http://localhost:8080/swagger/index.html

	// Start server
	return router.Run(":" + cfg.Port)
}

// Note: Make sure to implement logic for database operations and handle errors appropriately.
