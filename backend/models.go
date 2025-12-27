package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"unicode"
)

type CFProblem struct {
	ContestID int `json:"contestId"`
	Index string `json:"index"`
	Name string `json:"name"`
	Rating int `json:"rating"`
	Tags []string `json:"tags"`
}

type CFResponse struct {
	Status string `json:"status"`
	Result struct { Problems []CFProblem `json:"problems"`} `json:"result"`
}

func getProblems() ([]CFProblem, error) {
	resp, err := http.Get("https://codeforces.com/api/problemset.problems")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiData CFResponse
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

func saveProblemsToDB(problems []CFProblem) {
	tagMap := map[string]string{
		// foundation
		"implementation": "implementation",
		"brute force": "implementation",

		// ad-hoc
		"constructive algorithms": "ad-hoc",

		// sorting
		"sortings": "sortings",

		// two pointers
		"two pointers": "two pointers",

		// searching
		"binary search": "searching",
		"ternary search": "searching",
		"divide and conquer": "searching",

		// meet-in-the-middle
		"meet-in-the-middle": "meet-in-the-middle",

		// greedy
		"greedy": "greedy",

		// math + advanced math
		"math": "math",
		"number theory": "math",
		"combinatorics": "math",
		"matrices": "math",
		"probabilities": "math",
		"fft": "advanced math",
		"chinese remainder theorem": "advanced math",

		// geometry
		"geometry": "geometry",

		// graphs + advanced graphs
		"graphs": "graphs",
		"dfs and similar": "graphs",
		"shortest paths": "graphs",
		"dsu": "graphs",
		"flows": "advanced graphs",
		"graph matchings": "advanced graphs",
		"2-sat": "advanced graphs",

		// trees
		"trees": "trees",

		// strings + advanced strings
		"strings": "strings",
		"hashing": "strings",
		"string suffix structures": "advanced strings",

		// data structures
		"data structures": "data structures",
		"bitmasks": "data structures",

		// dp
		"dp": "dynamic programming",
	}
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
			topicSet["tree-dp"] = true
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

type Problem struct {
    ID string `json:"problem_id"`
	Name string `json:"name"`
    Rating int `json:"rating"`
    Tags []string `json:"tags"`
}