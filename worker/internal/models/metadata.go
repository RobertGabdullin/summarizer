package models

import "time"

type Metadata struct {
	Id           int
	Status       string
	Prompt       string
	Participants int
	CreatedAt    time.Time
}
