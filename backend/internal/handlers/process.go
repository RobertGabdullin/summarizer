package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/RobertGabdullin/summarizer/internal/models"
)

const (
	maxMemory    = 2000 << 20
	maxRAMMemory = 10 << 20
)

type DataSaver interface {
	SaveData(context.Context, *models.Record) error
}

type Process struct {
	dataSaver DataSaver
}

func (u *Process) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxMemory)
	if err := r.ParseMultipartForm(maxRAMMemory); err != nil {
		// to do
	}

	file, headers, err := r.FormFile("file")
	if err != nil {
		// to do
	}
	defer file.Close()

	processingType := r.FormValue("processingType")
	prompt := r.FormValue("summarizationPrompt")
	participantsStr := r.FormValue("participantsCount")

	participants, err := strconv.Atoi(participantsStr)
	if err != nil {
		// to do
	}

	metadata := &models.Metadata{
		Prompt:       prompt,
		Participants: participants,
	}
	record := &models.Record{
		Meeting:  file,
		Metadata: metadata,
	}

	err = u.dataSaver.SaveData(r.Context(), record)
	if err != nil {
		// to do
	}

	log.Printf("Запись успешно загружена!")
	log.Printf("Название файла: %s\nВид обработки: %s\nПромпт: %s\nКоличество участников: %d", headers.Filename, processingType, prompt, participants)
}
