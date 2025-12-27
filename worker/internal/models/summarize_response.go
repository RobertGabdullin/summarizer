package models

type SummarizeResponse struct {
	Choices []SummarizeChoice `json:"choices"`
}
