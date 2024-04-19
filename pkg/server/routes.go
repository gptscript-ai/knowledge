package server

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/acorn-io/z"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gptscript-ai/knowledge/pkg/db"
	"github.com/gptscript-ai/knowledge/pkg/types"
	"github.com/gptscript-ai/knowledge/pkg/types/defaults"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	golcdocloaders "github.com/hupe1980/golc/documentloader"
	golcschema "github.com/hupe1980/golc/schema"
	lcgodocloaders "github.com/tmc/langchaingo/documentloaders"
	lcgoschema "github.com/tmc/langchaingo/schema"
	"log/slog"
	"net/http"
	"path"
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

	// Create collection
	err := s.vs.CreateCollection(c, dataset.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		slog.Error("Failed to delete dataset (from DB)", "error", tx.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	err := s.vs.RemoveCollection(c, id)
	if err != nil {
		slog.Error("Failed to delete dataset (from vector store)", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	if query.TopK == nil {
		t := defaults.TopK
		query.TopK = &t
	}
	slog.Debug("Retrieving content from dataset", "dataset", id, "query", query)

	docs, err := s.vs.SimilaritySearch(c, query.Prompt, *query.TopK, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Retrieved documents", "count", len(docs), "docs", docs)
	c.JSON(http.StatusOK, gin.H{"results": docs})
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
		slog.Error("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate ID if none was provided
	if ingest.FileID == nil {
		fid, err := uuid.NewUUID()
		if err != nil {
			slog.Error("Failed to generate UUID", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ingest.FileID = z.Pointer(fid.String())
	}

	// decode content
	data, err := base64.StdEncoding.DecodeString(ingest.Content)
	if err != nil {
		slog.Error("Failed to decode content", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*
	 * Detect filetype
	 */
	var filetype string
	if ingest.Filename != nil {
		filetype = path.Ext(*ingest.Filename)
	}
	if filetype == "" {
		filetype = http.DetectContentType(data)
		if filetype == "application/octet-stream" { // fallback to mimetype
			filetype = mimetype.Detect(data).String()
		}
	}

	/*
	 * Set filename if not provided
	 */
	if ingest.Filename == nil {
		ingest.Filename = z.Pointer("<unnamed_document>")
	}

	slog.Debug("Loading data", "type", filetype, "filename", *ingest.Filename, "size", len(data))

	/*
	 * Load documents from the content
	 * For now, we're using documentloaders from both langchaingo and golc
	 * and translate them to our document schema.
	 */

	var lcgodocs []lcgoschema.Document
	var golcdocs []golcschema.Document

	switch filetype {
	case ".pdf", "application/pdf":
		r, err := golcdocloaders.NewPDF(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			slog.Error("Failed to create PDF loader", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		golcdocs, err = r.Load(c)
	case ".html", "text/html":
		lcgodocs, err = lcgodocloaders.NewHTML(bytes.NewReader(data)).Load(c)
	case ".md", ".txt", "text/plain", "text/markdown":
		lcgodocs, err = lcgodocloaders.NewText(bytes.NewReader(data)).Load(c)
	case ".csv", "text/csv":
		golcdocs, err = golcdocloaders.NewCSV(bytes.NewReader(data)).Load(c)
		if err != nil && errors.Is(err, csv.ErrBareQuote) {
			oerr := err
			err = nil
			var nerr error
			golcdocs, nerr = golcdocloaders.NewCSV(bytes.NewReader(data), func(o *golcdocloaders.CSVOptions) {
				o.LazyQuotes = true
			}).Load(c)
			if nerr != nil {
				err = errors.Join(oerr, nerr)
			}
		}
	case ".ipynb":
		golcdocs, err = golcdocloaders.NewNotebook(bytes.NewReader(data)).Load(c)
	default:
		slog.Error("Unsupported file type", "filename", *ingest.Filename)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type"})
		return
	}

	if err != nil {
		slog.Error("Failed to load document", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	docs := make([]vs.Document, len(lcgodocs)+len(golcdocs))
	for idx, doc := range lcgodocs {
		doc.Metadata["filename"] = *ingest.Filename
		docs[idx] = vs.Document{
			Metadata: doc.Metadata,
			Content:  doc.PageContent,
		}
	}

	for idx, doc := range golcdocs {
		doc.Metadata["filename"] = *ingest.Filename
		docs[idx] = vs.Document{
			Metadata: doc.Metadata,
			Content:  doc.PageContent,
		}

	}

	if len(docs) == 0 {
		slog.Error("No documents found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "no documents found"})
		return
	}

	// Add documents to VectorStore -> This generates the embeddings
	slog.Debug("Ingesting documents", "count", len(docs))
	docIDs, err := s.vs.AddDocuments(c, docs, id)
	if err != nil {
		slog.Error("Failed to add documents", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record file and documents in database
	dbDocs := make([]db.Document, len(docIDs))
	for idx, docID := range docIDs {
		dbDocs[idx] = db.Document{
			ID:      docID,
			FileID:  *ingest.FileID,
			Dataset: id,
		}
	}
	dbFile := db.File{
		ID:        *ingest.FileID,
		Dataset:   id,
		Documents: dbDocs,
	}

	tx := s.db.WithContext(c).Create(&dbFile)
	if tx.Error != nil {
		slog.Error("Failed to create file", "error", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": docIDs, "ingest": ingest})
}

// RemoveDocumentFromDataset removes a document from a dataset by ID. If the owning file context is now empty, the FileIndex is removed.
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

	// Find in Database
	var document db.Document
	tx := s.db.WithContext(c).First(&document, "id = ? AND dataset = ?", docID, id)
	if tx.Error != nil {
		slog.Error("Failed to find document", "error", tx.Error.Error(), "dataset", id, "document", docID)
		c.JSON(http.StatusNotFound, gin.H{"error": tx.Error.Error()})
		return
	}

	// Remove from VectorStore
	if err := s.vs.RemoveDocument(c, docID, id); err != nil {
		slog.Error("Failed to remove document", "error", err, "dataset", id, "document", docID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Remove from Database
	tx = s.db.WithContext(c).Delete(&document)
	if tx.Error != nil {
		slog.Error("Failed to delete document", "error", tx.Error.Error(), "dataset", id, "document", docID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	// Check if owning file should be removed
	var count int64
	tx = s.db.WithContext(c).Model(&db.Document{}).Where("file_id = ?", document.FileID).Count(&count)
	if tx.Error != nil {
		slog.Error("Failed to count documents", "error", tx.Error.Error(), "file", document.FileID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	if count == 0 {
		slog.Info("Removing file, because all associated documents are gone", "file", document.FileID)
		tx = s.db.WithContext(c).Delete(&db.File{}, "id = ?", document.FileID)
		if tx.Error != nil {
			slog.Error("Failed to delete file", "error", tx.Error.Error(), "file", document.FileID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
			return
		}
	}

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

	// Find file in database with associated documents
	var file db.File
	tx := s.db.WithContext(c).Preload("Documents").Where("id = ? AND dataset = ?", fileID, id).First(&file)
	if tx.Error != nil {
		slog.Error("Failed to find file", "error", tx.Error.Error(), "dataset", id, "file", fileID)
		c.JSON(http.StatusNotFound, gin.H{"error": tx.Error.Error()})
		return
	}

	// Remove owned documents from VectorStore and Database
	for _, doc := range file.Documents {
		if err := s.vs.RemoveDocument(c, doc.ID, id); err != nil {
			slog.Error("Failed to remove document", "error", err, "dataset", id, "document", doc.ID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		tx = s.db.WithContext(c).Delete(&doc)
		if tx.Error != nil {
			slog.Error("Failed to delete document", "error", tx.Error.Error(), "dataset", id, "document", doc.ID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
			return
		}
	}

	// Remove file DB
	tx = s.db.WithContext(c).Delete(&file)
	if tx.Error != nil {
		slog.Error("Failed to delete file", "error", tx.Error.Error(), "dataset", id, "file", fileID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

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

	// Get dataset with files and associated documents preloaded
	var dataset db.Dataset
	tx := s.db.WithContext(c).Preload("Files.Documents").First(&dataset, "id = ?", id)
	if tx.Error != nil {
		slog.Error("Failed to find dataset", "error", tx.Error.Error(), "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": tx.Error.Error()})
		return
	}

	slog.Info("Found dataset", "id", dataset.ID, "files", len(dataset.Files))

	c.JSON(http.StatusOK, dataset)
}
