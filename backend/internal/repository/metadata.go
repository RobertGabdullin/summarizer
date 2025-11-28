package repository

import (
	"context"

	"github.com/RobertGabdullin/summarizer/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	client *pgxpool.Pool
}

func NewPostgresClient(dsn string) (*PostgresClient, error) {
	pool, err := pgxpool.New(context.TODO(), dsn)
	if err != nil {
		return nil, err
	}

	client := &PostgresClient{
		client: pool,
	}

	return client, nil
}

func (c *PostgresClient) StoreMetadata(ctx context.Context, metadata *models.Metadata) (int, error) {
	sql := "INSERT INTO records () VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var id int
	err := c.client.QueryRow(ctx, sql).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
