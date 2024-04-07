package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	apiUrl = "https://api.openai.com/v1/chat/completions"
)

func isRateLimitError(response *resty.Response) bool {
	// Decodifique o JSON para uma estrutura de dados Go
	var responseMap map[string]interface{}
	if err := json.Unmarshal(response.Body(), &responseMap); err != nil {
		// Se ocorrer um erro ao decodificar o JSON, assuma que não é um erro de limite de taxa
		return false
	}

	// Verifique se o JSON contém um campo "error"
	if errorData, ok := responseMap["error"].(map[string]interface{}); ok {
		// Se o campo "error" estiver presente, verifique se contém informações sobre o limite de taxa
		if errorMessage, ok := errorData["message"].(string); ok {
			// Verifique se a mensagem de erro contém palavras-chave indicando um erro de limite de taxa
			return strings.Contains(errorMessage, "Rate limit reached") || strings.Contains(errorMessage, "rate_limit_exceeded")
		}
	}

	return false
}

var rateLimitCounter int

func sendRequest(query string) (string, error) {
	myquery := "Translate this text to Brazilian Portuguese reply only with the translation: " + query
	apiKey := "myKey"

	client := resty.New()
	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      "gpt-3.5-turbo",
			"messages":   []interface{}{map[string]interface{}{"role": "system", "content": myquery}},
			"max_tokens": 100,
		}).
		Post(apiUrl)

	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}

	if isRateLimitError(response) {
		rateLimitCounter++
		fmt.Println("Response Body:", string(response.Body()))
		fmt.Println(rateLimitCounter, "Erro de Rate limit reached na linha ", query, "\n", "Esperando ", 10, "s")
		time.Sleep(10 * time.Second)
		fmt.Println("Fim do periodo de espera")
		sendRequest(query)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response.Body(), &data); err != nil {
		return "", fmt.Errorf("error decoding JSON response: %v", err)
	}

	// Extract content with error handling
	choices, ok := data["choices"].([]interface{}) //???
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("invalid response structure: choices not found")
	}

	choice := choices[0].(map[string]interface{})

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response structure: message not found")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response structure: content not found")
	}

	return content, nil
}
