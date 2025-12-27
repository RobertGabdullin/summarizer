package models

const (
	systemRole = "system"
	userRole   = "user"
)

type SummarizeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewSystemMessage(promt string) SummarizeMessage {
	return SummarizeMessage{
		Role:    systemRole,
		Content: promt,
	}
}

func NewUserMessage(transcribation string) SummarizeMessage {
	return SummarizeMessage{
		Role:    userRole,
		Content: transcribation,
	}
}
