package main

import (
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

type AncestryMap map[string]map[string]int;

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
	Problem CFProblem `json:"problem"`
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

//builds the distance map
func BuildAncestryMap(nodes []Node, edges []Edge) AncestryMap {
	ancestry := make(AncestryMap)

	idToSlug := make(map[int]string)
	adjlist := make(map[int][]int)
	for _, node := range nodes {
		idToSlug[node.ID] = node.Slug
		ancestry[node.Slug] = make(map[string]int)
		ancestry[node.Slug][node.Slug] = 0
	}
	
	for _, edge := range edges {
		adjlist[edge.To] = append(adjlist[edge.To], edge.From)
	}

	//bfs through each node
	for _, node := range nodes {

		type pair struct {
			id int
			dist int
		}

		q := list.New()
		q.PushBack(pair{node.ID, 0})

		//bfs
		for q.Len() > 0 {
			elem := q.Front()
			cur := elem.Value.(pair)
			q.Remove(elem)

			curSlug := idToSlug[cur.id]
			ancestry[node.Slug][curSlug] = cur.dist
			
			for _, neighbor := range adjlist[cur.id] {
				neighborSlug := idToSlug[neighbor]
				if _, ok := ancestry[node.Slug][neighborSlug]; !ok {
					ancestry[node.Slug][neighborSlug] = cur.dist + 1
					q.PushBack(pair{neighbor, cur.dist+1})
				}
			}
		}
	}

	return ancestry
}

//calculates B
func GetBaseRating(rating int, attempts int) float64 {
	if attempts <= 1 {
		return float64(rating)
	}

	const k = 0.1

	modifier := 0.5 + 0.5*math.Exp(-k*float64(attempts-1))
	
	return float64(rating) * modifier
}

//calculates M given a B(j) and multipliers(j) for all j in the interval
func CalculateIntervalBin(solves []SolveAttributes) BinState {
	if len(solves) == 0 {
		return BinState{0, nil, nil}
	}

	var p float64 //max of c
	credits := make([]float64, len(solves)) //c array
	multiplier := make([]float64, len(solves)) //multipliers array

	for i, solve := range solves {
		credits[i] = solve.BaseRating * solve.Multiplier
		multiplier[i] = solve.Multiplier
		if credits[i] > p {
			p = credits[i]
		}
	}

	if p == 0 {
		return BinState{0, credits, multiplier}
	}

	var numerator, denominator float64

	const K = 1.5 //confidence constant

	for i, solve := range solves {
		weight := math.Pow((credits[i]/p), 3)

		numerator += credits[i] * weight
		denominator += solve.Multiplier * weight
	}
	denominator = math.Max(denominator, K)
	score := numerator/denominator

	return BinState{score, credits, multiplier}
}

//computed M(i, T) given T and the array of submissions at interval i. uses CalculateIntervalBin and GetBaseRating
func GetTopicIntervalState(targetTopic string, intervalSubmissions []Submission, ancestry AncestryMap) BinState {
	var attributes []SolveAttributes
	
	for _, submission := range intervalSubmissions {
		minDist := -1
		for _, topic := range submission.TopicSlugs {
			if dist, ok := ancestry[topic][targetTopic]; ok {
				if minDist == -1 || dist < minDist {
					minDist = dist
				}
			}
		}

		if minDist != -1 {
			base := GetBaseRating(submission.Rating, submission.Attempts)
			multipler := math.Pow(0.75, float64(minDist))
			attributes = append(attributes, SolveAttributes{base, multipler})
		}
	}
	return CalculateIntervalBin(attributes)
}

//calculates mastery score (cur and peak) given slice of interval scores
func CalculateMasteryScore(binScores []float64) MasteryResult {
	if len(binScores) == 0 {
		return MasteryResult{0, 0}
	}

	var p float64
	for _, score := range binScores {
		if score > p {
			p = score
		}
	}

	if p == 0 {
		return MasteryResult{0, 0}
	}

	const lambda = 0.05
	const K = 1.2

	var numerator float64
	var denominator float64

	for i, score := range binScores {
		timeWeight := math.Exp(-lambda * float64(i))
		qualityWeight := math.Pow(score/p, 3)

		totalWeight := timeWeight * qualityWeight

		numerator += score * totalWeight
		denominator += totalWeight
	}
	if denominator == 0 {
		return MasteryResult{0, 0}
	}
	return MasteryResult{numerator/math.Max(denominator, K), p}
}

//helper for GetBinnedSubmissions. returns index of bin given a time and int n
func GetAbsoluteBinIdx(t time.Time, n int) int {
    return int(t.Unix() / int64(n*86400))
}

//takes all submissions and an int n and groups them into n-day intervals
func GetBinnedSubmissions(submissions []Submission, n int) map[int][]Submission {
	binToSub := make(map[int][]Submission)
    
    for _, sub := range submissions {
        idx := GetAbsoluteBinIdx(sub.SolvedAt, n)
        binToSub[idx] = append(binToSub[idx], sub)
    }
    return binToSub
}

func IndexBinMap(binIdxToState map[int]BinState, minIdx int, currentBinIdx int) []float64 {
	var scoresForDecay []float64
	for i := currentBinIdx; i >= minIdx; i-- {
		if val, ok := binIdxToState[i]; ok {
			scoresForDecay = append(scoresForDecay, val.Score)
		} else {
			scoresForDecay = append(scoresForDecay, 0)
		}
	}
	return scoresForDecay
}

//returns a map, mapping each topic to its current mastery score and peak mastery score
func CalculateAllTopicMasteries(topics []string, submissions []Submission, ancestry AncestryMap, n int) (map[string]MasteryResult, map[string]map[int]BinState) {
	results := make(map[string]MasteryResult)
	allStates := make(map[string]map[int]BinState)
	binnedSubs := GetBinnedSubmissions(submissions, n)
	currentBinIdx := GetAbsoluteBinIdx(time.Now(), n)

	for _, topicSlug := range topics {
		binIdxToState := make(map[int]BinState)
		minIdx := -1

		for binIdx, intervalSubs := range binnedSubs {
			state := GetTopicIntervalState(topicSlug, intervalSubs, ancestry)
			if state.Score > 0{
				binIdxToState[binIdx] = state
				if minIdx == -1 || binIdx < minIdx {
					minIdx = binIdx
				}
			}
		}

		if len(binIdxToState) == 0 {
            results[topicSlug] = MasteryResult{0, 0}
            continue
        }

		scoresForDecay := IndexBinMap(binIdxToState, minIdx, currentBinIdx)

		cur := CalculateMasteryScore(scoresForDecay)

		results[topicSlug] = cur
		allStates[topicSlug] = binIdxToState
	}

	return results, allStates
}

//takes in the handle and other parameters. returns cur mastery score, peak mastery score, and problems solved/failed
func GetUserMastery(handle string, tagMap map[string]string, ancestry AncestryMap, n int) (map[string]MasteryResult, map[string]map[int]BinState, map[string]int, error) {
	url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s", handle)

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, nil, err
	}
	defer resp.Body.Close()

	var data CFUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil, nil, err
	}

	problemHistory := make(map[string][]CFSubmission)
	for _, submission := range data.Result {
		key := fmt.Sprintf("%d%s", submission.Problem.ContestID, submission.Problem.Index)
		problemHistory[key] = append(problemHistory[key], submission)
	}

	var processed []Submission
	problemsStatus := make(map[string]int)

	for id, subs := range problemHistory {
		var firstSolve *CFSubmission
		//problems are last to first so need to go from end to find first OK
		attempts := 1
		for i := len(subs) - 1; i >= 0; i-- {
			if subs[i].Verdict == "OK" {
				firstSolve = &subs[i];
				break
			}
			attempts++
		}

		if firstSolve == nil {
			problemsStatus[id] = -1 //-1 means incomplete
			continue
		}
		problemsStatus[id] = 1 //1 means complete
		
		var slugs []string
		var tree bool
		var dp bool
		for _, tag := range firstSolve.Problem.Tags {
			if topic, ok := tagMap[tag]; ok {
				slugs = append(slugs, topic)
			}
			if tag == "trees" {
				tree = true
			}
			if tag == "dp" {
				dp = true
			}
		}
		if tree && dp {
			slugs = append(slugs, "tree dp")
		}

		processed = append(processed, Submission{
			ID: id,
			Rating: firstSolve.Problem.Rating,
			Attempts: attempts,
			TopicSlugs: slugs,
			SolvedAt: time.Unix(firstSolve.CreationTimeSeconds, 0),
		})
	}
	topics := make([]string, 0, len(ancestry))
	for slug := range ancestry {
		topics = append(topics, slug)
	}
	masteryResults, binStates := CalculateAllTopicMasteries(topics, processed, ancestry, n)
	return masteryResults, binStates, problemsStatus, nil
}

func SyncUser(handle string, tagMap map[string]string, ancestry AncestryMap) error {
	const n = 14
	masteryResults, binStates, _, err := GetUserMastery(handle, tagMap, ancestry, n)
    if err != nil {
        return err
    }

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for topic, mastery := range masteryResults {
		_, err = tx.Exec(context.Background(), `
            INSERT INTO user_topic_stats (handle, topic_slug, mastery_score, peak_score, last_updated)
            VALUES ($1, $2, $3, $4, NOW())
            ON CONFLICT (handle, topic_slug) 
            DO UPDATE SET 
                mastery_score = EXCLUDED.mastery_score,
                peak_score = GREATEST(user_topic_stats.peak_score, EXCLUDED.peak_score),
                last_updated = NOW()`,
            handle, topic, mastery.Current, mastery.Peak)
		if err != nil {
			return err
		}
	}

	for topic, bin := range binStates {
		for binIdx, state := range bin {
			_, err = tx.Exec(context.Background(), `
                INSERT INTO user_interval_stats (handle, topic_slug, bin_idx, bin_score, credits, multipliers, last_updated)
                VALUES ($1, $2, $3, $4, $5, $6, NOW())
                ON CONFLICT (handle, topic_slug, bin_idx)
                DO UPDATE SET
                    bin_score = EXCLUDED.bin_score,
                    credits = EXCLUDED.credits,
                    multipliers = EXCLUDED.multipliers,
                    last_updated = NOW()`,
                handle, topic, binIdx, state.Score, state.Credits, state.Multipliers)
            if err != nil {
                return err
            }
		}
	}
	return tx.Commit(context.Background())
}

// func testing(handle string) {
// 	Nodes := make([]Node, 0, 100)
// 	Edges := make([]Edge, 0, 100)

// 	//getting nodes
// 	nodeRows, _ := conn.Query(context.Background(), "SELECT id, slug, display_name FROM topics")
//     defer nodeRows.Close()

// 	for nodeRows.Next() {
// 		var n Node
//         if err := nodeRows.Scan(&n.ID, &n.Slug, &n.DisplayName); err == nil {
//             Nodes = append(Nodes, n)
//         }
// 	}

// 	//getting edges
// 	edgeRows, _ := conn.Query(context.Background(), "SELECT parent_id, child_id FROM topic_dependencies")
//     defer edgeRows.Close()

//     for edgeRows.Next() {
//         var e Edge
//         if err := edgeRows.Scan(&e.From, &e.To); err == nil {
//             Edges = append(Edges, e)
//         }
//     }
// 	ancestry := BuildAncestryMap(Nodes, Edges)

// 	tagMap := map[string]string{
// 		// foundation
// 		"implementation": "implementation",
// 		"brute force": "implementation",

// 		// ad-hoc
// 		"constructive algorithms": "ad hoc",

// 		// sorting
// 		"sortings": "sortings",

// 		// two pointers
// 		"two pointers": "two pointers",

// 		// searching
// 		"binary search": "searching",
// 		"ternary search": "searching",
// 		"divide and conquer": "searching",

// 		// meet-in-the-middle
// 		"meet-in-the-middle": "meet in the middle",

// 		// greedy
// 		"greedy": "greedy",

// 		// math + advanced math
// 		"math": "math",
// 		"number theory": "math",
// 		"combinatorics": "math",
// 		"matrices": "math",
// 		"probabilities": "math",
// 		"fft": "advanced math",
// 		"chinese remainder theorem": "advanced math",

// 		// geometry
// 		"geometry": "geometry",

// 		// graphs + advanced graphs
// 		"graphs": "graphs",
// 		"dfs and similar": "graphs",
// 		"shortest paths": "graphs",
// 		"dsu": "graphs",
// 		"flows": "advanced graphs",
// 		"graph matchings": "advanced graphs",
// 		"2-sat": "advanced graphs",

// 		// trees
// 		"trees": "trees",

// 		// strings + advanced strings
// 		"strings": "strings",
// 		"hashing": "strings",
// 		"string suffix structures": "advanced strings",

// 		// data structures
// 		"data structures": "data structures",
// 		"bitmasks": "data structures",

// 		// dp
// 		"dp": "dynamic programming",
// 	}

// 	userMastery, _, _ := GetUserMastery(handle, tagMap, ancestry, 14)

// 	if len(userMastery) == 0 {
// 		fmt.Println("no data found")
// 	}
	
// 	for tag, mastery := range userMastery {
// 		fmt.Printf("[%s] Current: %.2f | Peak: %.2f\n", tag, mastery.Current, mastery.Peak)
// 	}
// } 