package mastery

import (
	"github.com/jackc/pgx/v5/pgxpool"
    "github.com/tanaydonde/cf-curriculum-planner/backend/internal/models"
)

type MasteryService struct {
    tagMap   map[string]string
    ancestry models.AncestryMap
	conn *pgxpool.Pool
}

func NewMasteryService(conn *pgxpool.Pool) *MasteryService {
    nodes, edges := models.GetGraph(conn)
    return &MasteryService{tagMap: GetTagMap(), ancestry: BuildAncestryMap(nodes, edges), conn: conn}
}

func (s *MasteryService) Sync(handle string) error {
    return syncUser(s.conn, handle, s.tagMap, s.ancestry)
}

func (s *MasteryService) RefreshAndGetAllStats(handle string) (map[string]MasteryResult, error) {
    return refreshAndGetAllStats(s.conn, handle, s.tagMap)
}

func (s *MasteryService) UpdateSubmission(handle string, problem ProblemSolveInput) error {
    return updateSubmissionFull(s.conn, handle, problem, s.tagMap, s.ancestry)
}

func (s *MasteryService) RecommendProblem(handle string, topic string, targetInc int, k int) ([]CFProblemOutput, error) {
    return recommendProblem(s.conn, handle, topic, targetInc, k)
}

func (s* MasteryService) RecommendDailyProblem(handle string) (CFProblemOutput, error) {
    return recommendDailyProblem(s.conn, handle)
}

func (s* MasteryService) GetLastKSolves(handle string, k int, status string) ([]CFSolveOutput, error) {
    return getLastKSolves(s.conn, handle, k, status)
}