package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AncestryMap map[string]map[string]int;

func GetGraph(conn *pgxpool.Pool) ([]Node, []Edge) {
    var nodes []Node
    var edges []Edge

    nRows, _ := conn.Query(context.Background(), "SELECT id, slug, display_name FROM topics")
    defer nRows.Close()
    for nRows.Next() {
        var n Node
        nRows.Scan(&n.ID, &n.Slug, &n.DisplayName)
        nodes = append(nodes, n)
    }

    eRows, _ := conn.Query(context.Background(), "SELECT parent_id, child_id FROM topic_dependencies")
    defer eRows.Close()
    for eRows.Next() {
        var e Edge
        eRows.Scan(&e.From, &e.To)
        edges = append(edges, e)
    }

    return nodes, edges
}

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

type Node struct {
    ID int `json:"id"`
    Slug string `json:"slug"`
    DisplayName string `json:"display_name"`
}

type Edge struct {
    From int `json:"from"`
    To int `json:"to"`
}
