package db

import (
    "context"
    "fmt"
    "os"
    "github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
    connStr := os.Getenv("DATABASE_URL")
    fmt.Println("CONNECTION STRING:", connStr != "")
    conn, err := pgxpool.New(context.Background(), connStr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "cannot connect to database: %v\n", err)
        os.Exit(1)
    }
    return conn
}