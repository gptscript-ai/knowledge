package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gptscript-ai/knowledge/pkg/db"
	"github.com/gptscript-ai/knowledge/pkg/docs"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore/chromem"
	cg "github.com/philippgille/chromem-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tmc/langchaingo/vectorstores"
	"log/slog"
	"net/http"
)

type Config struct {
	ServerURL, Port, APIBase string
}

type Server struct {
	db *db.DB
	vs vectorstores.VectorStore
}

func NewServer(db *db.DB) *Server {
	return &Server{db: db}
}

// Start starts the server with the given configuration.
func (s *Server) Start(ctx context.Context, cfg Config) error {
	router := gin.Default()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	// Database migration
	if err := s.db.AutoMigrate(); err != nil {
		return err
	}

	// API routes
	docs.SwaggerInfo.BasePath = cfg.APIBase
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	router.GET("/swagger/*any", swaggerHandler) // http://localhost:8080/swagger/index.html

	// @title Knowledge API
	// @version 1
	// @description This is the Knowledge API server for GPTStudio.
	// @contact.name Acorn Labs Inc.
	v1 := router.Group(cfg.APIBase)
	{
		// Swagger >>>
		v1.GET("/docs/*any", swaggerHandler)
		v1.GET("/docs", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, cfg.APIBase+"/docs/index.html")
		})
		// <<< Swagger

		v1Datasets := v1.Group("/datasets")
		{
			v1Datasets.GET("/", s.ListDatasets)
			v1Datasets.GET("/:id", s.GetDataset)
			v1Datasets.POST("/create", s.CreateDataset)
			v1Datasets.DELETE("/:id", s.DeleteDataset)
			v1Datasets.POST("/:id/ingest", s.IngestIntoDataset)
			v1Datasets.POST("/:id/query", s.QueryDataset)
			v1Datasets.POST("/:id/retrieve", s.RetrieveFromDataset)
			v1Datasets.DELETE("/:id/documents/:doc_id", s.RemoveDocumentFromDataset)
			v1Datasets.DELETE("/:id/files/:file_id", s.RemoveFileFromDataset)
		}
	}

	vsdb, err := cg.NewPersistentDB("vector.db", false)
	if err != nil {
		return err
	}

	normalized := true
	ef := cg.NewEmbeddingFuncOpenAICompat(
		"http://localhost:8080/v1",
		"sk-foo",
		"text-embedding-ada-002",
		&normalized,
	)

	col, err := vsdb.GetOrCreateCollection("default", nil, ef)
	if err != nil {
		return err
	}

	vs, err := chromem.New(chromem.WithDB(vsdb), chromem.WithEmbeddingFunc(ef), chromem.WithCollection(col))
	if err != nil {
		return err
	}
	s.vs = vs

	// Start server
	return router.Run(":" + cfg.Port)
}

// Note: Make sure to implement logic for database operations and handle errors appropriately.
