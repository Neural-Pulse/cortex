package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neural-pulse/cortex/backend/internal/app/elasticsearch"
	"github.com/neural-pulse/cortex/backend/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	logger, err := logging.ConfigureLogger()
	if err != nil {
		log.Fatalf("Failed to configure logger: %v", err)
	}

	logger.Info("Aplicação iniciada")

	esClient, err := elasticsearch.NewElasticsearchClient([]string{"https://localhost:9200"}, "elastic", "RbXM5XOGW-PpTl9HDonA", "http_ca.crt")
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}

	r := gin.Default()

	r.POST("/configurar-banco", func(c *gin.Context) {
		dataBaseConfig(c, esClient, logger)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func dataBaseConfig(c *gin.Context, esClient *elasticsearch.ElasticsearchClient, logger *zap.Logger) {
	document := map[string]interface{}{"example": "data"}

	// Corrigido para chamar IndexDocument diretamente no esClient
	err := esClient.IndexDocument(c.Request.Context(), "your_index", "", document)
	if err != nil {
		logger.Error("Failed to index document", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to index document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document indexed successfully"})
}
