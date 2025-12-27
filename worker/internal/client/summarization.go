package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/RobertGabdullin/summarizer/internal/models"
)

type Summarizer struct {
	baseUrl string
	client  *http.Client
}

func NewSummarizer(baseUrl string) *Summarizer {
	return &Summarizer{
		baseUrl: baseUrl,
		client:  http.DefaultClient,
	}
}

func (c *Summarizer) GetSummarization(promt, transcribation string) (string, error) {
	summarizeBody := models.NewSummarizeRequest(promt, transcribation)

	jsonSummarizeBody, err := json.Marshal(summarizeBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseUrl+"/v1/chat/completions", bytes.NewBuffer(jsonSummarizeBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer no-key-needed")

	log.Println("Выполняем запрос к модели саммаризации")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	summarizationJson, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var summarizeResp models.SummarizeResponse
	err = json.Unmarshal(summarizationJson, &summarizeResp)
	if err != nil {
		return "", fmt.Errorf("Не получилось конвертировать json в структуру: %w", err)
	}

	if len(summarizeResp.Choices) == 0 {
		return "", fmt.Errorf("Некорретный саммарайз")
	}

	return summarizeResp.Choices[0].Message.Content, nil
}
