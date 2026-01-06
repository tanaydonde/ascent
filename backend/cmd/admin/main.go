package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tanaydonde/cf-curriculum-planner/backend/internal/db"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../../../app.env")
    if err != nil {
        fmt.Fprintf(os.Stderr, "can't load .env")
    }

	conn := db.Connect()
	defer conn.Close()

	fmt.Println("successfully connected to the database")

	script, err := os.ReadFile("../../internal/db/init.sql")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read SQL file: %v\n", err)
		os.Exit(1)
	}

	_, err = conn.Exec(context.Background(), string(script))
	if err != nil {
		fmt.Fprintf(os.Stderr, "SQL execution failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("starting database seeding")
	db.FillTables(conn)

	fmt.Println("seeding complete, database now ready")
}