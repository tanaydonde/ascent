package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tanaydonde/cf-curriculum-planner/backend/internal/mastery"
	"github.com/tanaydonde/cf-curriculum-planner/backend/internal/models"
)

func FillTables(conn *pgxpool.Pool) {
	tagMap := mastery.GetTagMap()
	saveProblemsToDB(tagMap, conn)
	createTopics(tagMap, conn)
	createRoadMap(conn)
}

func saveProblemsToDB(tagMap map[string]string, conn *pgxpool.Pool) {
	problems, _ := getProblems()
	for _, p := range problems {
		if p.Rating == 0 {
			continue
		}

		if cyrillic(p.Name) {
			continue
		}

		topicSet := make(map[string]bool)
		hasDP := false
		hasTrees := false
		for _, tag := range p.Tags {
			if topic, ok := tagMap[tag]; ok {
				topicSet[topic] = true
				if topic == "dynamic programming" {
					hasDP = true
				}
				if topic == "trees" {
					hasTrees = true
				}
			}
		}
		if hasDP && hasTrees {
			topicSet["tree dp"] = true
		}

		filtered := make([]string, 0, len(topicSet))
		for topic := range topicSet {
			filtered = append(filtered, topic)
		}
		
		if len(filtered) == 0 {
			continue
		}

		problemID := fmt.Sprintf("%d%s", p.ContestID, p.Index)
		
		query := `
			INSERT INTO problems (problem_id, name, rating, tags)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (problem_id) DO UPDATE
			SET rating = EXCLUDED.rating, tags = EXCLUDED.tags;
		`
		_, err := conn.Exec(context.Background(), query, problemID, p.Name, p.Rating, filtered)
		if err != nil {
			fmt.Printf("could not save problem %s: %v\n", problemID, err)
		}
	}
	//fmt.Println("total problem count:", count)
	fmt.Println("all rated problems saved successfully")
}

func createTopics(tagMap map[string]string, conn *pgxpool.Pool) {
	uniqueTopics := make(map[string]bool)
	for _, topicSlug := range tagMap {
		uniqueTopics[topicSlug] = true
	}

	for slug := range uniqueTopics {
		query := `
			INSERT INTO topics (slug, display_name)
			VALUES ($1, $2)
			ON CONFLICT (slug) DO UPDATE
			SET display_name = EXCLUDED.display_name
		`
		_, err := conn.Exec(context.Background(), query, slug, getDisplayName(slug))
		if err != nil {
			fmt.Printf("could not save topic %s: %v\n", slug, err)
		}
	}
	// for tree dp
	slug := "tree dp"
	query := `
		INSERT INTO topics (slug, display_name)
		VALUES ($1, $2)
		ON CONFLICT (slug) DO UPDATE
		SET display_name = EXCLUDED.display_name
	`
	_, err := conn.Exec(context.Background(), query, slug, getDisplayName(slug))
	if err != nil {
		fmt.Printf("could not save topic %s: %v\n", slug, err)
	}
}

func createRoadMap(conn *pgxpool.Pool) {
	// 19 edges in total
	linkTopics("implementation", "ad hoc", conn)
	linkTopics("implementation", "sortings", conn)
	linkTopics("implementation", "data structures", conn)
	linkTopics("implementation", "greedy", conn)
	linkTopics("implementation", "math", conn)
	linkTopics("implementation", "strings", conn)

	linkTopics("sortings", "two pointers", conn)
	linkTopics("sortings", "searching", conn)

	linkTopics("data structures", "searching", conn)
	linkTopics("data structures", "graphs", conn)

	linkTopics("greedy", "dynamic programming", conn)

	linkTopics("math", "advanced math", conn)
	linkTopics("math", "geometry", conn)

	linkTopics("strings", "advanced strings", conn)

	linkTopics("searching", "meet in the middle", conn)

	linkTopics("dynamic programming", "tree dp", conn)

	linkTopics("trees", "tree dp", conn)

	linkTopics("graphs", "advanced graphs", conn)
	linkTopics("graphs", "trees", conn)
}

func getProblems() ([]models.CFProblem, error) {
	resp, err := http.Get("https://codeforces.com/api/problemset.problems")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiData models.CFResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiData); err != nil {
		return nil, err
	}

	return apiData.Result.Problems, nil
}

func cyrillic(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Cyrillic, r) {
			return true
		}
	}
	return false
}

func getDisplayName(topic string) string {
	switch topic {
	case "tree dp":
		return "Tree DP"
	case "dynamic programming":
		return "DP"
	}

	words := strings.Fields(topic)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func linkTopics(parent string, child string, conn *pgxpool.Pool) {
	query := `
		INSERT INTO topic_dependencies (parent_id, child_id)
		SELECT p.id, c.id FROM topics p, topics c WHERE p.slug = $1 AND c.slug = $2
		ON CONFLICT (parent_id, child_id) DO NOTHING
	`
	_, err := conn.Exec(context.Background(), query, parent, child)
    if err != nil {
        fmt.Printf("error linking %s -> %s: %v\n", parent, child, err)
        return
    }
}