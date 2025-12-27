package main

import (
	"log"

	"github.com/RobertGabdullin/summarizer/internal/config"
	"github.com/RobertGabdullin/summarizer/internal/launcher"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Println("Не вышло получить конфиг")
		return
	}

	launcer := launcher.New(cfg)
	err = launcer.Start()
	if err != nil {
		log.Println("Не вышло стартануть лаунчер")
		return
	}
}
