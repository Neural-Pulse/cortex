package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"go.uber.org/zap"
)

// ConfigureLogger configura e retorna uma instância do logger zap.
func ConfigureLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// SendLogToElasticsearch envia um log para o Elasticsearch.
func SendLogToElasticsearch(es *elasticsearch.Client, message string, logger *zap.Logger) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   message,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(logEntry); err != nil {
		logger.Error("Erro ao codificar logEntry para JSON", zap.Error(err))
		return
	}

	req := esapi.IndexRequest{
		Index:   "logs-index",
		Body:    &buf,
		Refresh: "true",
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		logger.Error("Erro ao enviar log para Elasticsearch", zap.Error(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error("Resposta de erro ao enviar log para Elasticsearch", zap.String("resposta", res.String()))
	}
}
