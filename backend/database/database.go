package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	log.Println("Postgres Connection String:", os.Getenv("POSTGRES_CONN"))
	pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN"))
	if err != nil {
		log.Printf("Unable to create database pool: %v", err)
		return nil, err
	}
	return pool, nil
}
