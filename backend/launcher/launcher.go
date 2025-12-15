package launcher

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RobertGabdullin/summarizer/internal/config"
	"github.com/RobertGabdullin/summarizer/internal/handlers"
	"github.com/RobertGabdullin/summarizer/internal/repository"
	"github.com/RobertGabdullin/summarizer/internal/service"
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
	postgresRepository, err := getPostgresRepository(l.cfg.Postgres)
	if err != nil {
		return fmt.Errorf("Не получилось получить Postgresql репозиторий: %w", err)
	}
	s3Repository, err := getS3Repository(l.cfg.S3)
	if err != nil {
		return fmt.Errorf("Не получилось получить S3 репозиторий: %w", err)
	}

	dataService := service.NewUserData(postgresRepository, s3Repository)

	processHandler := handlers.NewProcess(dataService)

	http.Handle("/api/process", processHandler)
	http.Handle("/", handlers.NewStatic())

	log.Println("Listening on http://localhost:1717")
	http.ListenAndServe(":1717", nil)

	return nil
}

func getPostgresRepository(cfg *config.PostgresConfig) (*repository.PostgresClient, error) {
	return repository.NewPostgresClient(cfg.User, cfg.Password, cfg.Host, cfg.DBName, cfg.Port)
}

func getS3Repository(cfg *config.S3Config) (*repository.S3Client, error) {
	return repository.NewS3Client(cfg.Region, cfg.Bucket, cfg.AccessKeyID, cfg.SecretAccessKey, cfg.Endpoint)
}
