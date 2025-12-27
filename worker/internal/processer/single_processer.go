package processer

import (
	"fmt"
	"io"

	"github.com/RobertGabdullin/summarizer/internal/client"
)

type SingleProcesser struct {
	summarizationClient  *client.Summarizer
	transcribationClient *client.Transcribator
}

func NewSingleProcesser(summarizationClient *client.Summarizer, transcribationClient *client.Transcribator) *SingleProcesser {
	return &SingleProcesser{
		summarizationClient:  summarizationClient,
		transcribationClient: transcribationClient,
	}
}

func (p *SingleProcesser) ProcessMeeting(meeting io.ReadCloser, prompt string) (string, string, error) {

	transcribation, err := p.transcribationClient.GetTranscribation(meeting)
	if err != nil {
		return "", "", fmt.Errorf("Не вышло получить транскрибацию: %w", err)
	}

	summarization, err := p.summarizationClient.GetSummarization(prompt, transcribation)
	if err != nil {
		return "", "", fmt.Errorf("Не вышло получить саммаризацию: %w", err)
	}

	return transcribation, summarization, nil
}
