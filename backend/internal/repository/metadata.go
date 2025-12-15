package repository

import (
	"context"
	"fmt"

	"github.com/RobertGabdullin/summarizer/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	client *pgxpool.Pool
}

func NewPostgresClient(user, password, host, dbname string, port int) (*PostgresClient, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname)

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
	sql := "INSERT INTO records (prompt, participants) VALUES ($1, $2) RETURNING id"
	var id int
	err := c.client.QueryRow(ctx, sql, metadata.Prompt, metadata.Participants).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
