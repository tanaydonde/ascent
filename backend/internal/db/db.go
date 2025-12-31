package db

import (
    "context"
    "fmt"
    "os"
    "github.com/jackc/pgx/v5"
)

func Connect() *pgx.Conn {
    connStr := "postgres://donde783985@localhost:5432/cf_planner"
    conn, err := pgx.Connect(context.Background(), connStr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "cannot connect to database: %v\n", err)
        os.Exit(1)
    }
    return conn
}