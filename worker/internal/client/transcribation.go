package client

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type Transcribator struct {
	baseUrl string
	client  *http.Client
}

func NewTranscrcibator(baseUrl string) *Transcribator {
	return &Transcribator{
		baseUrl: baseUrl,
		client:  http.DefaultClient,
	}
}

func (c *Transcribator) GetTranscribation(meetingFile io.ReadCloser) (string, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer writer.Close()
		defer meetingFile.Close()

		part, err := writer.CreateFormFile("file", "some file")
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		if _, err := io.Copy(part, meetingFile); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	req, err := http.NewRequest(http.MethodPost, c.baseUrl+"/transcrip", pr)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	log.Println("Выполняем запрос к модели транскрибации")
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	transcribation, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(transcribation), nil
}
