package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getProblemsByTopic(w http.ResponseWriter, r *http.Request) {
	topic := chi.URLParam(r, "topic")

	query := `SELECT problem_id, name, rating, tags FROM problems WHERE $1 = ANY(tags) LIMIT 50`
	rows, err := conn.Query(context.Background(), query, topic)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var results []Problem
	for rows.Next() {
		var p Problem
		if err := rows.Scan(&p.ID, &p.Name, &p.Rating, &p.Tags); err != nil {
			continue
		}
		results = append(results, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}