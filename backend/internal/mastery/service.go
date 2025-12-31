package mastery

import (
	"github.com/jackc/pgx/v5"
    "github.com/tanaydonde/cf-curriculum-planner/backend/internal/models"
)

type MasteryService struct {
    tagMap   map[string]string
    ancestry models.AncestryMap
	conn *pgx.Conn
}

func NewMasteryService(conn *pgx.Conn) *MasteryService {
    nodes, edges := models.GetGraph(conn)
    return &MasteryService{tagMap: GetTagMap(), ancestry: BuildAncestryMap(nodes, edges), conn: conn}
}

func (s *MasteryService) Sync(handle string) error {
    return syncUser(s.conn, handle, s.tagMap, s.ancestry)
}

//to do later
func (s *MasteryService) GetCurrentMastery(handle string, topic string) float64 {
    return 0.0 
}

//to do later
func (s *MasteryService) GetPeakMastery(handle string, topic string) float64 {
    return 0.0 
}