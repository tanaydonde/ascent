package mastery

import (
	"time"
	"github.com/tanaydonde/cf-curriculum-planner/backend/internal/models"
)

type SolveAttributes struct {
	BaseRating float64
	Multiplier float64
}

type Submission struct {
	ID string
	Rating int
	Attempts int
	TopicSlugs []string
	SolvedAt time.Time
}

type MasteryResult struct {
	Current float64
	Peak    float64
}

type CFSubmission struct {
	Verdict string `json:"verdict"`
	Problem models.CFProblem `json:"problem"`
	CreationTimeSeconds int64 `json:"creationTimeSeconds"`
}

type CFUserResponse struct {
	Status string `json:"status"`
	Result []CFSubmission `json:"result"`
}

type BinState struct {
	Score float64
    Credits []float64
    Multipliers []float64
}
