package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gptscript-ai/knowledge/pkg/db"
	"github.com/gptscript-ai/knowledge/pkg/types"
	"github.com/gptscript-ai/knowledge/pkg/types/defaults"
	"log/slog"
	"net/http"
)

// CreateDataset creates a new dataset.
// @Summary Create a new dataset
// @Description Create a new dataset
// @Tags datasets
// @Accept json
// @Produce json
// @Param dataset body types.Dataset true "Dataset object"
// @Router /datasets/create [post]
func (s *Server) CreateDataset(c *gin.Context) {
	var dataset types.Dataset
	if err := c.ShouldBindJSON(&dataset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if dataset.EmbedDimension == nil || *dataset.EmbedDimension <= 0 {
		f := defaults.EmbeddingDimension
		dataset.EmbedDimension = &f
	}

	// Create dataset
	slog.Info("Creating dataset", "id", dataset.ID)
	tx := s.db.WithContext(c).Create(&dataset)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, dataset)
}

// DeleteDataset deletes a dataset by ID.
// @Summary Delete a dataset
// @Description Delete a dataset by ID
// @Tags datasets
// @Produce json
// @Param id path string true "Dataset ID"
// @Router /datasets/{id} [delete]
func (s *Server) DeleteDataset(c *gin.Context) {
	id := c.Param("id")
	slog.Info("Deleting dataset", "id", id)

	tx := s.db.WithContext(c).Delete(&types.Dataset{}, "id = ?", id)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// QueryDataset queries a dataset by ID.
// @Summary Query a dataset
// @Description Query a dataset by ID
// @Tags datasets
// @Accept json
// @Produce json
// @Param id path string true "Dataset ID"
// @Router /datasets/{id}/query [post]
func (s *Server) QueryDataset(c *gin.Context) {
	id := c.Param("id")
	slog.Info("Querying dataset", "id", id)

	var query types.Query
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Query", "query", query)

	// validate
	v := validator.New()
	if err := v.Struct(query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	d := s.db.WithContext(c)

	ds, err := db.GetDataset(d, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("dataset not found: %v", err), "id": id})
		return
	}

	slog.Info("Dataset", "id", ds.ID)

	// TODO: DB query logic here
	c.JSON(http.StatusOK, gin.H{"id": id, "query": query})
}

// RetrieveFromDataset retrieves content from a dataset by ID.
// @Summary Retrieve content from a dataset
// @Description Retrieve content from a dataset by ID
// @Tags datasets
// @Accept json
// @Produce json
// @Param id path string true "Dataset ID"
// @Router /datasets/{id}/retrieve [post]
func (s *Server) RetrieveFromDataset(c *gin.Context) {
	id := c.Param("id")
	slog.Info("Retrieving content from dataset", "dataset", id)

	var query types.Query
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: DB query logic here
	c.JSON(http.StatusOK, gin.H{"id": id, "query": query})
}

// IngestIntoDataset ingests content into a dataset by ID.
// @Summary Ingest content into a dataset
// @Description Ingest content into a dataset by ID
// @Tags datasets
// @Accept json
// @Produce json
// @Param id path string true "Dataset ID"
// @Router /datasets/{id}/ingest [post]
func (s *Server) IngestIntoDataset(c *gin.Context) {
	id := c.Param("id")
	slog.Info("Ingesting content into dataset", "dataset", id)

	var ingest types.Ingest
	if err := c.ShouldBindJSON(&ingest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: DB ingest logic here
	c.JSON(http.StatusOK, gin.H{"id": id, "ingest": ingest})
}

// RemoveDocumentFromDataset removes a document from a dataset by ID.
// @Summary Remove a document from a dataset
// @Description Remove a document from a dataset by ID
// @Tags datasets
// @Accept json
// @Produce json
// @Param id path string true "Dataset ID"
// @Param doc_id path string true "Document ID"
// @Router /datasets/{id}/documents/{doc_id} [delete]
func (s *Server) RemoveDocumentFromDataset(c *gin.Context) {
	id := c.Param("id")
	docID := c.Param("doc_id")
	slog.Info("Removing document from dataset", "dataset", id, "document", docID)
	// TODO: DB remove logic here
	c.JSON(http.StatusOK, gin.H{"id": id, "doc_id": docID})
}

// RemoveFileFromDataset removes a file from a dataset by ID.
// @Summary Remove a file from a dataset
// @Description Remove a file from a dataset by ID
// @Tags datasets
// @Accept json
// @Produce json
// @Param id path string true "Dataset ID"
// @Param file_id path string true "File ID"
// @Router /datasets/{id}/files/{file_id} [delete]
func (s *Server) RemoveFileFromDataset(c *gin.Context) {
	id := c.Param("id")
	fileID := c.Param("file_id")
	slog.Info("Removing file from dataset", "dataset", id, "file", fileID)
	// TODO: DB remove logic here
	c.JSON(http.StatusOK, gin.H{"id": id, "file_id": fileID})
}

// ListDatasets lists all datasets.
// @Summary List all datasets
// @Description List all datasets
// @Tags datasets
// @Produce json
// @Router /datasets [get]
func (s *Server) ListDatasets(c *gin.Context) {
	slog.Info("Listing datasets")
	tx := s.db.WithContext(c).Find(&[]types.Dataset{})
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return

	}

	var datasets []types.Dataset
	if err := tx.Scan(&datasets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, datasets)
}

// GetDataset gets a dataset by ID.
// @Summary Get a dataset
// @Description Get a dataset by ID
// @Tags datasets
// @Produce json
// @Param id path string true "Dataset ID"
// @Router /datasets/{id} [get]
func (s *Server) GetDataset(c *gin.Context) {
	id := c.Param("id")
	slog.Info("Getting dataset", "id", id)
	dataset := &types.Dataset{}
	tx := s.db.WithContext(c).First(dataset, "id = ?", id)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, dataset)
}
