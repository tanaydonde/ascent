package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tanaydonde/cf-curriculum-planner/backend/internal/db"
)

func main() {

	err := godotenv.Load("../../../app.env")
    if err != nil {
        fmt.Fprintf(os.Stderr, "can't load .env")
    }

	conn := db.Connect()

	for i := 0; i < 15; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := conn.Ping(ctx)
		cancel()

		if err == nil {
			break
		}

		if i == 14 {
			log.Fatalf("database not reachable after retries: %v", err)
		}

		sleep := time.Duration(min(200*(1<<i), 60000)) * time.Millisecond
		time.Sleep(sleep)
	}


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