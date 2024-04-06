package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	elasticsearch_utils "github.com/neural-pulse/cortex-backend/utils"

	elasticsearch "github.com/elastic/go-elasticsearch/v8" // Importe a versão v8 correspondente
	_ "github.com/go-sql-driver/mysql"
)

// Configuração fixa do cliente Elasticsearch
var es *elasticsearch.Client

func main() {
	// Configurar o cliente Elasticsearch
	esConfig := elasticsearch_utils.ElasticsearchConfig{
		Addresses: []string{"https://localhost:9200"},
		Username:  "elastic",
		Password:  "HQLANIxpisY2jGhsMQ*F",
	}

	var err error
	es, err = elasticsearch_utils.NewElasticsearchClient(esConfig)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.POST("/configurar-banco", configurarBanco)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func configurarBanco(c *gin.Context) {
	var req struct {
		DSN        string `json:"dsn"`
		DBType     string `json:"dbType"`
		ESURL      string `json:"esUrl"`
		ESUser     string `json:"esUser"`
		ESPass     string `json:"esPass"`
		CACertPath string `json:"caCertPath"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	esConfig := elasticsearch_utils.ElasticsearchConfig{
		Addresses: []string{"https://localhost:9200"},
		Username:  "elastic",
		Password:  "HQLANIxpisY2jGhsMQ*F",
	}

	// Conectar ao banco de dados
	db, err := sql.Open(req.DBType, req.DSN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	// Configurar o cliente Elasticsearch
	es, err := elasticsearch_utils.NewElasticsearchClient(esConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Listar todos os bancos de dados disponíveis
	databases, err := elasticsearch_utils.ListDatabases(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Lista de bancos de dados excluídos
	excludedDatabases := []string{"sys", "information_schema", "mysql", "performance_schema", "information_schema"}

	// Iterar sobre todos os bancos de dados
	for _, database := range databases {
		// Verificar se o banco de dados está na lista de excluídos
		if contains(excludedDatabases, database) {
			continue // Ignorar bancos de dados excluídos
		}

		fmt.Printf("Banco de dados: %s\n", database)

		// Listar tabelas e colunas para o banco de dados atual
		if err := elasticsearch_utils.ListTablesAndColumns(db, database, es); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuração concluída com sucesso"})
}

// Função auxiliar para verificar se um slice contém um determinado valor
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
