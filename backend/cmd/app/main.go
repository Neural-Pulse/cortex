package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/neural-pulse/cortex/backend/internal/app/database"
	"github.com/neural-pulse/cortex/backend/internal/app/elasticsearch"
	"github.com/neural-pulse/cortex/backend/internal/app/schema"
	"github.com/neural-pulse/cortex/backend/pkg/logging"
	"go.uber.org/zap"
)

type ConfigRequest struct {
	DSN    string `json:"dsn"`
	DBType string `json:"dbType"`
}

func encrypt(data []byte, passphrase string) (string, error) {
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(ciphertext), nil
}

func main() {
	logger, err := logging.ConfigureLogger()
	if err != nil {
		log.Fatalf("Failed to configure logger: %v", err)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	mariaDB, err := database.SetupMariaDB(user, password, host, dbName)
	if err != nil {
		logger.Error("Error setup database", zap.Error(err))
	}

	logger.Info("App Starts")

	esClient, err := elasticsearch.NewElasticsearchClient([]string{"https://localhost:9200"}, "elastic", "RbXM5XOGW-PpTl9HDonA", "http_ca.crt")
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}

	r := gin.Default()

	r.POST("/get-db-data", func(c *gin.Context) {
		dataBaseInfoConfig(c, esClient, logger, mariaDB)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func dataBaseInfoConfig(c *gin.Context, esClient *elasticsearch.ElasticsearchClient, logger *zap.Logger, mariaDB *database.MariaDB) {
	var req ConfigRequest
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON"})
		return
	}

	dbInstance := database.NewDatabaseFactory(req.DBType)
	db, err := dbInstance.SetupDatabase(req.DSN)
	if err != nil {
		logger.Error("Failed to setup database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to setup database"})
		return
	}

	schemas, err := schema.ListSchemas(db)
	if err != nil {
		logger.Error("Failed to list schemas", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list schemas"})
		return
	}

	for _, schemaName := range schemas {
		err := schema.ListTablesAndColumns(db, schemaName, esClient) // Ajuste conforme a implementação real
		if err != nil {
			logger.Error("Failed to list tables and columns for schema: "+schemaName, zap.Error(err))
			// Decida se quer parar o processo aqui ou apenas logar o erro e continuar
		}
	}
	config := &database.DatabaseConfig {
		DSN: req.DSN,
		DBType: req.DBType,
	}
	err = database.SaveDatabaseConfig(mariaDB.DB, config)
	if err != nil {
		logger.Error("Failed to save DSN and DBType to MariaDB", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save DSN and DBType to MariaDB"})
		return
	}


	c.JSON(http.StatusOK, gin.H{"message": "Database info indexed successfully"})
}
