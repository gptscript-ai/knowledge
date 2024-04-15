package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gptscript-ai/knowledge/pkg/db"
	"github.com/gptscript-ai/knowledge/pkg/types"
	"log/slog"
	"net/http"
)

// CreateDataset creates a new dataset.
// @Summary Create a new dataset
// @Description Create a new dataset
// @Tags datasets
// @Accept json
// @Produce json
// @Router /datasets/create [post]
func CreateDataset(c *gin.Context) {
	var dataset db.Dataset
	if err := c.ShouldBindJSON(&dataset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Creating dataset: %v", dataset)
	// TODO: DB insert logic here
	c.JSON(http.StatusOK, dataset)
}

// DeleteDataset deletes a dataset by ID.
// @Summary Delete a dataset
// @Description Delete a dataset by ID
// @Tags datasets
// @Produce json
// @Param id path string true "Dataset ID"
// @Router /datasets/{id} [delete]
func DeleteDataset(c *gin.Context) {
	id := c.Param("id")
	slog.Info("Deleting dataset: %s", id)
	// TODO: DB delete logic here
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
func QueryDataset(c *gin.Context) {
	id := c.Param("id")
	slog.Info("Querying dataset: %s", id)

	var query types.Query
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Query: %v", query)
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
func RetrieveFromDataset(c *gin.Context) {
	id := c.Param("id")
	slog.Info("Retrieving content from dataset: %s", id)

	var query types.Query
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Query: %v", query)

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
func IngestIntoDataset(c *gin.Context) {
	id := c.Param("id")
	slog.Info("Ingesting content into dataset: %s", id)

	var ingest types.Ingest
	if err := c.ShouldBindJSON(&ingest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Ingest: %v", ingest)

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
func RemoveDocumentFromDataset(c *gin.Context) {
	id := c.Param("id")
	docID := c.Param("doc_id")
	slog.Info("Removing document from dataset: %s, %s", id, docID)
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
func RemoveFileFromDataset(c *gin.Context) {
	id := c.Param("id")
	fileID := c.Param("file_id")
	slog.Info("Removing file from dataset: %s, %s", id, fileID)
	// TODO: DB remove logic here
	c.JSON(http.StatusOK, gin.H{"id": id, "file_id": fileID})
}

// ListDatasets lists all datasets.
// @Summary List all datasets
// @Description List all datasets
// @Tags datasets
// @Produce json
// @Router /datasets [get]
func ListDatasets(c *gin.Context) {
	slog.Info("Listing datasets")
	// TODO: DB list logic here
	c.JSON(http.StatusOK, gin.H{})
}

// GetDataset gets a dataset by ID.
// @Summary Get a dataset
// @Description Get a dataset by ID
// @Tags datasets
// @Produce json
// @Param id path string true "Dataset ID"
// @Router /datasets/{id} [get]
func GetDataset(c *gin.Context) {
	id := c.Param("id")
	slog.Info("Getting dataset: %s", id)
	// TODO: DB get logic here
	c.JSON(http.StatusOK, gin.H{"id": id})
}
