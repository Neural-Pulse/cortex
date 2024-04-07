package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var es *elasticsearch.Client

func main() {
	caCert, err := ioutil.ReadFile("http_ca.crt")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	esConfig := elasticsearch.Config{
		Addresses: []string{"https://localhost:9200"},
		Username:  "elastic",
		Password:  "soVw2YhEfHb1*H8Tf7K_",
		Transport: httpClient.Transport,
	}

	es, err = elasticsearch.NewClient(esConfig)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Ping(es.Ping.WithContext(context.Background()))
	if err != nil {
		log.Fatalf("Error pinging Elasticsearch: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Elasticsearch ping failed: %s", res.String())
	}

	log.Println("Successfully pinged Elasticsearch")

	r := gin.Default()

	r.POST("/configurar-banco", dataBaseConfig)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func dataBaseConfig(c *gin.Context) {
	if es == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Elasticsearch client not initialized"})
		return
	}
	var req struct {
		DSN    string `json:"dsn"`
		DBType string `json:"dbType"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := sql.Open(req.DBType, req.DSN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	// Call fetchAndIndexAllMetadata to fetch and index all metadata
	err = fetchAndIndexAllMetadata(db, es, "metadata_index")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Metadata indexed successfully"})
}

func fetchAndIndexAllMetadata(db *sql.DB, es *elasticsearch.Client, indexName string) error {
	// Fetch all databases (schemas)
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var databaseName string
		if err := rows.Scan(&databaseName); err != nil {
			return err
		}

		// Fetch all tables for the current database
		// Note: Directly concatenating the databaseName into the query
		tableRows, err := db.Query("SHOW TABLES FROM " + databaseName)
		if err != nil {
			return err
		}
		defer tableRows.Close()

		for tableRows.Next() {
			var tableName string
			if err := tableRows.Scan(&tableName); err != nil {
				return err
			}

			// Fetch all columns for the current table
			// Note: Directly concatenating the databaseName and tableName into the query
			columnRows, err := db.Query("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '" + databaseName + "' AND TABLE_NAME = '" + tableName + "'")
			if err != nil {
				return err
			}
			defer columnRows.Close()

			for columnRows.Next() {
				var columnName string
				if err := columnRows.Scan(&columnName); err != nil {
					return err
				}

				// Index the metadata into Elasticsearch
				metadata := map[string]interface{}{
					"database_name": databaseName,
					"table_name":    tableName,
					"column_name":   columnName,
				}
				err = indexMetadataIntoElasticsearch(es, indexName, []map[string]interface{}{metadata})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func indexMetadataIntoElasticsearch(es *elasticsearch.Client, indexName string, metadata []map[string]interface{}) error {
	for _, item := range metadata {
		// Convert the item map to JSON
		jsonStr, err := json.Marshal(item)
		if err != nil {
			return err
		}

		// Create a request body from the JSON string
		req := esapi.IndexRequest{
			Index:      indexName,
			DocumentID: "", // Let Elasticsearch auto-generate the ID
			Body:       strings.NewReader(string(jsonStr)),
			Refresh:    "true",
		}

		// Perform the request with the client.
		res, err := req.Do(context.Background(), es)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		// Check for errors in the response.
		if res.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				return err
			}
			// Print the response status and error information.
			return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
		}
	}
	return nil
}
