package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Postgres *PostgresConfig
	S3       *S3Config
}

func New() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	var postgresCfg PostgresConfig
	err = env.Parse(&postgresCfg)
	if err != nil {
		return nil, err
	}

	var s3Cfg S3Config
	err = env.Parse(&s3Cfg)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Postgres: &postgresCfg,
		S3:       &s3Cfg,
	}

	return &cfg, nil
}

type PostgresConfig struct {
	Host     string `env:"SUMMARIZER_POSTGRES_HOSTNAME" envDefault:"localhost"`
	Port     int    `env:"SUMMARIZER_POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"SUMMARIZER_POSTGRES_USERNAME" envDefault:"postgres"`
	Password string `env:"SUMMARIZER_POSTGRES_PASSWORD" envDefault:"postgres"`
	DBName   string `env:"SUMMARIZER_POSTGRES_DBNAME" envDefault:"summarizer"`
}

type S3Config struct {
	Region          string `env:"SUMMARIZER_AWS_REGION" envDefault:"localhost"`
	AccessKeyID     string `env:"SUMMARIZER_AWS_ACCESS_KEY_ID" envDefault:"minio"`
	SecretAccessKey string `env:"SUMMARIZER_AWS_SECRET_ACCESS_KEY" envDefault:"minio"`
	Endpoint        string `env:"SUMMARIZER_S3_ENDPOINT" envDefault:"http://localhost:9000"`
	Bucket          string `env:"SUMMARIZER_S3_BUCKET" envDefault:"meetings"`
}
