package elasticsearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockESClient é um mock para o ESClient.
type MockESClient struct {
	mock.Mock
}

// Index é um método mock que simula a indexação de um documento.
func (m *MockESClient) Index(ctx context.Context, index, documentID string, body interface{}) error {
	args := m.Called(ctx, index, documentID, body)
	return args.Error(0)
}

// TestIndexDocument testa a funcionalidade de indexação de documentos.
func TestIndexDocument(t *testing.T) {
	mockESClient := new(MockESClient)
	esClient := &ElasticsearchClient{client: mockESClient} // Instancia diretamente com o mock.

	// Configura o mock para esperar a chamada específica e retornar nil (nenhum erro).
	mockESClient.On("Index", mock.Anything, "test-index", "test-id", mock.AnythingOfType("map[string]interface {}")).Return(nil)

	// Chama a função que queremos testar, passando um contexto.
	err := esClient.IndexDocument(context.Background(), "test-index", "test-id", map[string]interface{}{"field": "value"})

	// Verifica se a função Index do mock foi chamada conforme esperado.
	mockESClient.AssertExpectations(t)

	// Verifica se não houve erro na indexação do documento.
	assert.NoError(t, err)
}
