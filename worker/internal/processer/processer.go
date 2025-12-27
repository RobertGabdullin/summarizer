package processer

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/RobertGabdullin/summarizer/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Querier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type MetadataRepository interface {
	LockNextPending(ctx context.Context, querier Querier) (*models.Metadata, error)
	SetFailed(ctx context.Context, id int, querier Querier) error
	SetSuccess(ctx context.Context, id int, transcribation string, summarization string, querier Querier) error
}

type MeetingProcesser interface {
	ProcessMeeting(meeting io.ReadCloser, prompt string) (string, string, error)
}

type MeetingRetriever interface {
	GetMeeting(id string) (io.ReadCloser, error)
}

type Processer struct {
	client             *pgxpool.Pool
	metadataRepository MetadataRepository
	meetingRetriever   MeetingRetriever
	meetingProcesser   MeetingProcesser
}

func NewMainProccesser(
	user, password, host, dbname string,
	port int,
	metadataRepository MetadataRepository,
	meetingRetriever MeetingRetriever,
	meetingProcesser MeetingProcesser,
) (*Processer, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname)

	pool, err := pgxpool.New(context.TODO(), dsn)
	if err != nil {
		return nil, err
	}

	processer := &Processer{
		client:             pool,
		metadataRepository: metadataRepository,
		meetingRetriever:   meetingRetriever,
		meetingProcesser:   meetingProcesser,
	}

	return processer, nil
}

func (p *Processer) RetrieveAndProcess(ctx context.Context) error {
	log.Println("Обработка началась")
	tx, err := p.client.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Не получилось начать транзакцию: %w", err)
	}
	defer tx.Rollback(ctx)

	log.Println("Успешно начали транзакцию")

	metadata, err := p.metadataRepository.LockNextPending(ctx, tx)
	if err != nil {
		return err
	}

	log.Println("Успешно получили метаданные")

	meeting, err := p.meetingRetriever.GetMeeting(fmt.Sprint(metadata.Id))
	if err != nil {
		return fmt.Errorf("Не получилось извлечь файл по id: %w", err)
	}

	log.Println("Успешно получили запись из S3")

	transcribation, summarization, err := p.meetingProcesser.ProcessMeeting(meeting, metadata.Prompt)
	if err != nil {
		return fmt.Errorf("Не вышло обработать запись: %w", err)
	}

	log.Println("Успешно обработали запись")

	err = os.WriteFile("transcribation.txt", []byte(transcribation), 0644)
	if err != nil {
		return fmt.Errorf("Не удалось записать транскрибацию в файл: %w", err)
	}

	err = os.WriteFile("summarization.txt", []byte(summarization), 0644)
	if err != nil {
		return fmt.Errorf("Не удалось записать саммари в файл: %w", err)
	}

	err = p.metadataRepository.SetSuccess(ctx, metadata.Id, transcribation, summarization, tx)
	if err != nil {
		return fmt.Errorf("Не получилось выставить статус: %w", err)
	}

	log.Println("Успешно загрузили транскрибацию и саммаризацию")

	tx.Commit(ctx)
	return nil
}
