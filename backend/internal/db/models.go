package db

type Problem struct {
    ID string `json:"problem_id"`
	Name string `json:"name"`
    Rating int `json:"rating"`
    Tags []string `json:"tags"`
}