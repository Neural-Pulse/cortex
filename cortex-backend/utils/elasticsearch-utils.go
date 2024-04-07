package elasticsearch_utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/go-sql-driver/mysql"
)

type ElasticsearchConfig struct {
	Addresses []string
	Username  string
	Password  string
}

// Estrutura para os campos adicionais de um documento
type AdditionalFields struct {
	Description        string   `json:"description"`
	DataClassification string   `json:"data_classification"`
	Tags               []string `json:"tags"`
	Health             string   `json:"health"`
}

// Função para preencher os campos adicionais de um documento
func FillAdditionalFields(schema, tablename, field string) AdditionalFields {
	// Aqui você pode adicionar a lógica para preencher os campos adicionais com base nos parâmetros fornecidos
	// Por exemplo:
	description := ""
	dataClassification := ""
	tags := []string{} // exemplo de tags
	health := ""       // exemplo de saúde do campo

	return AdditionalFields{
		Description:        description,
		DataClassification: dataClassification,
		Tags:               tags,
		Health:             health,
	}
}

func ListDatabases(db *sql.DB) ([]string, error) {
	var databases []string

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		databases = append(databases, dbName)
	}

	return databases, nil
}

func ListTablesAndColumns(db *sql.DB, schema string, es *elasticsearch.Client) error {
	tablesQuery := fmt.Sprintf("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s'", schema)
	rows, err := db.Query(tablesQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tablename string
		if err := rows.Scan(&tablename); err != nil {
			return err
		}
		fmt.Println("  Tabela:", tablename)

		columnsQuery := fmt.Sprintf("DESCRIBE %s.%s", schema, tablename)
		cols, err := db.Query(columnsQuery) // Aqui é onde o erro pode estar ocorrendo
		if err != nil {
			return err
		}
		defer cols.Close()

		for cols.Next() {
			var (
				field string
				typ   string
				null  string
				key   string
				def   sql.NullString
				extra string
			)
			if err := cols.Scan(&field, &typ, &null, &key, &def, &extra); err != nil {
				return err
			}

			// Preencher campos adicionais
			additionalFields := FillAdditionalFields(schema, tablename, field)

			// Adicionar campos adicionais aos dados do campo
			columnData := map[string]interface{}{
				"db":                  schema,
				"name":                field,
				"type":                typ,
				"allow_null":          null,
				"key":                 key,
				"default":             def,
				"extra":               extra,
				"description":         additionalFields.Description,
				"data_classification": additionalFields.DataClassification,
				"tags":                additionalFields.Tags,
				"health":              additionalFields.Health,
			}
			data, err := json.Marshal(columnData)
			if err != nil {
				return err
			}

			// Enviar dados para o Elasticsearch
			req := esapi.IndexRequest{
				Index:      "catalogo-teste",
				DocumentID: fmt.Sprintf("%s_%s_%s", schema, tablename, field),
				Body:       strings.NewReader(string(data)),
				Refresh:    "true",
			}

			res, err := req.Do(context.Background(), es)
			if err != nil {
				return err
			}
			defer res.Body.Close()
		}
	}
	return nil
}
