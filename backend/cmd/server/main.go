package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tanaydonde/cf-curriculum-planner/backend/internal/api"
	"github.com/tanaydonde/cf-curriculum-planner/backend/internal/db"
	"github.com/tanaydonde/cf-curriculum-planner/backend/internal/mastery"
)

func main() {
	conn := db.Connect()
	defer conn.Close(context.Background())

	service := mastery.NewMasteryService(conn)

	h := &api.Handler{Conn: conn, Service: service,}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // Your frontend URLs
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/problems/{topic}", h.GetProblemsByTopic)
		r.Get("/graph", h.GetGraphHandler)
		r.Post("/sync/{handle}", h.SyncUserHandler)
	})

	port := ":8080"
	fmt.Printf("Server starting on %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}