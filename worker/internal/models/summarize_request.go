package models

const (
	model       = "llama-3-8b"
	temperature = 0.7
	maxTokens   = 5000
)

type SummarizeRequest struct {
	Model       string             `json:"model"`
	Messages    []SummarizeMessage `json:"messages"`
	Temperature float32            `json:"temperature"`
	MaxTokens   int                `json:"max_tokens"`
}

func NewSummarizeRequest(promt, transcribation string) SummarizeRequest {
	messages := []SummarizeMessage{
		NewSystemMessage(promt),
		NewUserMessage(transcribation),
	}

	return SummarizeRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}
}
