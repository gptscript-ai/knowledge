package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gptscript-ai/knowledge/pkg/db"
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
	// DB insert logic here
	c.JSON(http.StatusOK, dataset)
}
