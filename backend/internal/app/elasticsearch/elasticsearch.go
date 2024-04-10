package elasticsearch

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// ESClient define as operações suportadas pelo cliente Elasticsearch.
type ESClient interface {
	Index(ctx context.Context, index, documentID string, body interface{}) error
}

// ElasticsearchAdapter é um adaptador que implementa a interface ESClient
// para o cliente Elasticsearch real.
type ElasticsearchAdapter struct {
	client *elasticsearch.Client
}

func (adapter *ElasticsearchAdapter) Index(ctx context.Context, index, documentID string, body interface{}) error {
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

	res, err := req.Do(ctx, adapter.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("erro ao indexar documento: %s", res.String())
	}

	return nil
}

// ElasticsearchClient é uma abstração sobre o cliente Elasticsearch real.
type ElasticsearchClient struct {
	client ESClient
}

// NewElasticsearchClient cria uma nova instância do cliente Elasticsearch.
func NewElasticsearchClient(addresses []string, username, password, caCertPath string) (*ElasticsearchClient, error) {
	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{RootCAs: caCertPool}
	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
		Transport: httpClient.Transport,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	adapter := &ElasticsearchAdapter{client: es}
	return &ElasticsearchClient{client: adapter}, nil
}

// IndexDocument indexa um documento no Elasticsearch usando o adaptador.
func (e *ElasticsearchClient) IndexDocument(ctx context.Context, index, documentID string, body interface{}) error {
	return e.client.Index(ctx, index, documentID, body)
}
