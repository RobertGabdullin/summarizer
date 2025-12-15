package main

import (
	"log"

	"github.com/RobertGabdullin/summarizer/internal/config"
	"github.com/RobertGabdullin/summarizer/launcher"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Println("Не вышло получить конфиг: %w", err)
		return
	}

	launcher := launcher.New(cfg)
	err = launcher.Start()
	if err != nil {
		log.Println("Не вышло стартануть лаунчер: %w", err)
	}
}
