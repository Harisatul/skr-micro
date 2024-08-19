package leaderboard

type gradeMessage struct {
	Id       string `json:"id"`
	Token    string `json:"token"`
	TryoutID string `json:"tryout_id"`
	Score    int64  `json:"score"`
	Version  int64  `json:"version"`
}

type leaderboardReq struct {
	ToId string `json:"tid"`
	Size string `json:"size"`
	Page string `json:"page"`
}

// Define a struct to hold the query results
type ParticipantScore struct {
	ParticipantID   string
	ParticipantName string
	Score           int
}
