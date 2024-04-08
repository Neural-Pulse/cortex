package elasticsearchservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

func SearchEmptyDescriptions(es *elasticsearch.Client) ([]string, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"description.keyword": "",
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %s", err)
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("metadata_index3"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	var columnNames []string
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		columnName := source["column_name"].(string)
		columnNames = append(columnNames, columnName)
	}

	return columnNames, nil
}

// StartService inicia o serviço de consulta ao Elasticsearch e envia prompts para a API LLM.
func StartService(esConfig elasticsearch.Config) {
	go func() {
		// Cria um novo cliente Elasticsearch com a configuração fornecida
		es, err := elasticsearch.NewClient(esConfig)
		if err != nil {
			log.Fatalf("Erro ao criar o cliente Elasticsearch: %s", err)
		}

		// Log para confirmar que o serviço iniciou
		log.Println("Iniciando o serviço de consulta ao Elasticsearch...")

		// Loop infinito que executa a função de busca a cada 3 minutos
		for {
			columns, err := SearchEmptyDescriptions(es)
			if err != nil {
				log.Printf("Erro ao buscar colunas com descrição vazia: %s", err)
				continue
			}

			if len(columns) > 0 {
				fmt.Println("Colunas com descrição vazia:")
				prompt := fmt.Sprintf("Por favor, forneça uma descrição para cada uma das seguintes colunas de bancos de dados %s, seguindo o formato 'nome_da_coluna:descrição, nome_da_coluna:descrição, ...'. Certifique-se de separar cada par 'nome_da_coluna:descrição' por vírgula (,).", strings.Join(columns, ", "))
				response := sendPromptToLLM("AIzaSyBovLANQbWmMZTqph7PKv9CPvXD5jT8ohE", prompt)
				processLLMResponse(es, response, columns)
			} else {
				fmt.Println("Nenhuma coluna com descrição vazia encontrada.")
			}

			// Aguarda 3 minutos antes de executar a próxima iteração
			time.Sleep(3 * time.Minute)
		}
	}()
}

// sendPromptToLLM envia um prompt para a API LLM e imprime a resposta.
func sendPromptToLLM(apiKey string, prompt string) string {
	// Prepara o corpo da requisição
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": prompt,
					},
				},
			},
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Erro ao preparar o corpo da requisição: %s", err)
		return ""
	}

	// Cria a requisição HTTP
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Erro ao criar a requisição: %s", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")

	// Envia a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao enviar a requisição: %s", err)
		return ""
	}
	defer resp.Body.Close()

	// Lê a resposta
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler a resposta: %s", err)
		return ""
	}

	log.Printf("Resposta da API: %s", string(responseBody))
	return string(responseBody)
}

func processLLMResponse(es *elasticsearch.Client, response string, columnNames []string) {
	var apiResponse map[string]interface{}
	err := json.Unmarshal([]byte(response), &apiResponse)
	if err != nil {
		log.Printf("Erro ao fazer parse da resposta JSON: %s", err)
		return
	}

	candidates, ok := apiResponse["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		log.Println("Nenhuma resposta de candidato encontrada na resposta da API LLM")
		return
	}

	firstCandidate, ok := candidates[0].(map[string]interface{})
	if !ok {
		log.Println("Candidato inválido encontrado na resposta da API LLM")
		return
	}

	content, ok := firstCandidate["content"].(map[string]interface{})
	if !ok {
		log.Println("Conteúdo inválido encontrado na resposta da API LLM")
		return
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		log.Println("Nenhuma parte de texto encontrada na resposta da API LLM")
		return
	}

	text := parts[0].(map[string]interface{})["text"].(string)

	re := regexp.MustCompile(`([^\s:]+):\s*([^,]+)`)
	matches := re.FindAllStringSubmatch(text, -1)

	descriptions := make(map[string]string)
	for _, match := range matches {
		if len(match) == 3 {
			fieldName := strings.TrimSpace(match[1])
			fieldDescription := strings.TrimSpace(match[2])
			descriptions[fieldName] = fieldDescription
		}
	}

	fmt.Println("Descrições:")
	for _, columnName := range columnNames {
		if description, ok := descriptions[columnName]; ok {
			fmt.Printf("%s: %s\n", columnName, description)
			updateQuery := map[string]interface{}{
				"query": map[string]interface{}{
					"term": map[string]interface{}{
						"column_name.keyword": columnName,
					},
				},
				"script": map[string]interface{}{
					"source": "ctx._source.description = params.description",
					"params": map[string]interface{}{
						"description": description,
					},
				},
			}

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(updateQuery); err != nil {
				log.Printf("Erro ao codificar a consulta de atualização: %s", err)
				continue
			}

			res, err := es.UpdateByQuery(
				[]string{"metadata_index3"},
				es.UpdateByQuery.WithBody(&buf),
				es.UpdateByQuery.WithContext(context.Background()),
				es.UpdateByQuery.WithPretty(),
			)
			if err != nil {
				log.Printf("Erro ao atualizar o documento no Elasticsearch: %s", err)
				continue
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("Erro na resposta do Elasticsearch: %s", res.String())
			} else {
				log.Printf("Documento atualizado com sucesso para a coluna %s", columnName)
			}
		} else {
			fmt.Printf("%s: Descrição não encontrada\n", columnName)

		}
	}
}
