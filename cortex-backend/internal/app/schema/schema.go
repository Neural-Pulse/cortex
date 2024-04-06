package schema

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8" // Importe a versão 8 do cliente Elasticsearch
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Schema representa um esquema de banco de dados.
type Schema struct {
	Name string
}

// Contains verifica se o nome do esquema está contido em uma lista de esquemas excluídos.
func (s Schema) Contains(excludedDatabases []string) bool {
	for _, db := range excludedDatabases {
		if s.Name == db {
			return true
		}
	}
	return false
}

// ListTablesAndColumns lista as tabelas e colunas para o esquema atual no banco de dados.
func (s Schema) ListTablesAndColumns(db *sql.DB, es *elasticsearch.Client) error {
	// Implemente a lógica para listar tabelas e colunas para o esquema atual aqui.
	return nil
}

func ListSchemas(db *sql.DB) ([]string, error) {
	var schemas []string

	rows, err := db.Query("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, err
		}
		schemas = append(schemas, schema)
	}

	return schemas, nil
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
		cols, err := db.Query(columnsQuery)
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

			// Montar dados do campo para Elasticsearch
			columnData := map[string]interface{}{
				"db":                  schema,
				"table":               tablename,
				"field":               field,
				"type":                typ,
				"allow_null":          null,
				"key":                 key,
				"default":             def,
				"extra":               extra,
				"description":         "", // Preencher com a descrição apropriada
				"data-classification": "", // Preencher com a classificação apropriada
				"tags":                "", // Preencher com as tags apropriadas
				"health":              "", // Preencher com o status de saúde apropriado
			}

			// Indexar dados no Elasticsearch
			data, err := json.Marshal(columnData)
			if err != nil {
				return err
			}

			req := esapi.IndexRequest{
				Index:      "catalogo",
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
