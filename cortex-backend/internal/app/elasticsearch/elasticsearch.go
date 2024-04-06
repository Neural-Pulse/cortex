package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// ElasticsearchClient é uma estrutura para interagir com o Elasticsearch.
type ElasticsearchClient struct {
	client *elasticsearch.Client
}

// NewElasticsearchClient cria uma nova instância do cliente Elasticsearch.
func NewElasticsearchClient(addresses []string, username, password string) (*ElasticsearchClient, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ElasticsearchClient{client: es}, nil
}

// IndexDocument indexa um documento no Elasticsearch.
func (e *ElasticsearchClient) IndexDocument(index, documentID string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: documentID,
		Body:       strings.NewReader(string(data)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("erro ao indexar documento: %s", res.String())
	}

	return nil
}

// UpdateDocument atualiza um documento existente no Elasticsearch.
func (e *ElasticsearchClient) UpdateDocument(index, documentID string, updateBody map[string]interface{}) error {
	data, err := json.Marshal(updateBody)
	if err != nil {
		return err
	}

	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: documentID,
		Body:       strings.NewReader(string(data)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("erro ao atualizar documento: %s", res.String())
	}

	return nil
}
