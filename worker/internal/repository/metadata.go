package repository

import (
	"context"
	"fmt"

	"github.com/RobertGabdullin/summarizer/internal/models"
	"github.com/RobertGabdullin/summarizer/internal/processer"
)

const (
	successStatus = "SUCCESS"
	failedStatus  = "FAILED"
)

type PostgresClient struct{}

func (c *PostgresClient) LockNextPending(ctx context.Context, querier processer.Querier) (*models.Metadata, error) {
	sql := "SELECT id, prompt, participants, status, created_at " +
		"FROM records WHERE status = 'PENDING' " +
		"ORDER BY created_at LIMIT 1 FOR UPDATE SKIP LOCKED"
	row := querier.QueryRow(ctx, sql)

	var metadata models.Metadata
	err := row.Scan(&metadata.Id, &metadata.Prompt, &metadata.Participants, &metadata.Status, &metadata.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("Не вышло извлечь запись: %w", err)
	}

	return &metadata, nil
}

func (c *PostgresClient) SetFailed(ctx context.Context, id int, querier processer.Querier) error {
	sql := "UPDATE records SET status = $1 WHERE id = $2"
	_, err := querier.Exec(ctx, sql, failedStatus, id)
	if err != nil {
		return fmt.Errorf("Не вышло изменить статус: %w", err)
	}

	return nil
}

func (c *PostgresClient) SetSuccess(ctx context.Context, id int, transcribation, summarization string, querier processer.Querier) error {
	sql := "UPDATE records SET transcribation = $1, summarization = $2, status = $3 WHERE id = $4"
	_, err := querier.Exec(ctx, sql, transcribation, summarization, successStatus, id)
	if err != nil {
		return fmt.Errorf("Не вышло обновить статус на успех: %w", err)
	}

	return nil
}
