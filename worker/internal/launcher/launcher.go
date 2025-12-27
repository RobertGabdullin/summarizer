package launcher

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/RobertGabdullin/summarizer/internal/client"
	"github.com/RobertGabdullin/summarizer/internal/config"
	"github.com/RobertGabdullin/summarizer/internal/processer"
	"github.com/RobertGabdullin/summarizer/internal/repository"
)

const (
	goroutineCount = 1
)

type Launcher struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Launcher {
	return &Launcher{
		cfg: cfg,
	}
}

func (l *Launcher) Start() error {

	summarizeClient := client.NewSummarizer(l.cfg.Summarizer.BaseUrl)
	transcribationClient := client.NewTranscrcibator(l.cfg.Transcribation.BaseUrl)

	singleProccesser := processer.NewSingleProcesser(summarizeClient, transcribationClient)

	metadataRepository := &repository.PostgresClient{}
	meetingRepository, err := repository.NewS3Client(
		l.cfg.S3.Region,
		l.cfg.S3.Bucket,
		l.cfg.S3.AccessKeyID,
		l.cfg.S3.SecretAccessKey,
		l.cfg.S3.Endpoint,
	)
	if err != nil {
		return fmt.Errorf("Не вышло создать S3 клиент: %w", err)
	}

	mainProcesser, err := processer.NewMainProccesser(
		l.cfg.Postgres.User,
		l.cfg.Postgres.Password,
		l.cfg.Postgres.Host,
		l.cfg.Postgres.DBName,
		l.cfg.Postgres.Port,
		metadataRepository,
		meetingRepository,
		singleProccesser,
	)
	if err != nil {
		return fmt.Errorf("Не удалось создать процессера: %w", err)
	}

	wg := &sync.WaitGroup{}

	for range goroutineCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				err := mainProcesser.RetrieveAndProcess(context.TODO())
				if err != nil {
					log.Println(err.Error())
				}
				time.Sleep(15)
			}
		}()
	}

	wg.Wait()
	return nil
}
