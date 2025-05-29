package dto

import "time"

// CreateMatchRequest 매칭 생성 요청
type CreateMatchRequest struct {
	GameID     string `json:"gameId"`
	MaxPlayers int    `json:"maxPlayers"`
}

// CreateMatchResponse 매칭 생성 응답
type CreateMatchResponse struct {
	MatchID    string `json:"matchId"`
	GameID     string `json:"gameId"`
	HostID     string `json:"hostId"`
	MaxPlayers int    `json:"maxPlayers"`
	Message    string `json:"message"`
}

// InviteFriendRequest 친구 초대 요청
type InviteFriendRequest struct {
	MatchID   string   `json:"matchId"`
	FriendIds []string `json:"friendIds"`
}

// InviteFriendResponse 친구 초대 응답
type InviteFriendResponse struct {
	MatchID    string   `json:"matchId"`
	InvitedIds []string `json:"invitedIds"`
	FailedIds  []string `json:"failedIds"`
	Message    string   `json:"message"`
}

// MatchInvitation 매칭 초대 알림
type MatchInvitation struct {
	MatchID   string `json:"matchId"`
	GameID    string `json:"gameId"`
	HostID    string `json:"hostId"`
	HostName  string `json:"hostName,omitempty"`
	ExpiresAt int64  `json:"expiresAt"`
	Message   string `json:"message"`
}

// InviteResponseRequest 초대 응답 요청
type InviteResponseRequest struct {
	MatchID  string `json:"matchId"`
	Response string `json:"response"` // "accept" or "decline"
}

// InviteResponseResponse 초대 응답 결과
type InviteResponseResponse struct {
	MatchID  string `json:"matchId"`
	UserID   string `json:"userId"`
	Response string `json:"response"`
	Message  string `json:"message"`
}

// MatchInfo 매칭 정보
type MatchInfo struct {
	MatchID    string        `json:"matchId"`
	GameID     string        `json:"gameId"`
	HostID     string        `json:"hostId"`
	Status     string        `json:"status"` // "waiting", "ready", "starting", "playing"
	Players    []MatchPlayer `json:"players"`
	MaxPlayers int           `json:"maxPlayers"`
	CreatedAt  time.Time     `json:"createdAt"`
	StartedAt  *time.Time    `json:"startedAt,omitempty"`
}

// MatchPlayer 매칭 플레이어 정보
type MatchPlayer struct {
	UserID   string     `json:"userId"`
	Status   string     `json:"status"` // "invited", "joined", "ready"
	JoinedAt *time.Time `json:"joinedAt,omitempty"`
}

// StartMatchRequest 매칭 시작 요청
type StartMatchRequest struct {
	MatchID string `json:"matchId"`
}

// StartMatchResponse 매칭 시작 응답
type StartMatchResponse struct {
	MatchID string `json:"matchId"`
	GameID  string `json:"gameId"`
	Teams   []Team `json:"teams"`
	Message string `json:"message"`
}

// Team 팀 정보
type Team struct {
	ID      int      `json:"id"`
	Players []string `json:"players"`
}

// GameInfo 게임 정보
type GameInfo struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"`
	PlayerCount int                    `json:"playerCount"`
	Teams       []Team                 `json:"teams"`
	Players     []MatchPlayer          `json:"players"`
	CreatedAt   time.Time              `json:"createdAt"`
	StartedAt   *time.Time             `json:"startedAt,omitempty"`
	EndedAt     *time.Time             `json:"endedAt,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// GameEventRequest 게임 이벤트 요청
type GameEventRequest struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// GameStartedResponse 게임 시작 응답
type GameStartedResponse struct {
	GameID string `json:"gameId"`
	Status string `json:"status"`
}

// PlayerDisconnectedResponse 플레이어 연결 해제 응답
type PlayerDisconnectedResponse struct {
	PlayerID string `json:"playerId"`
	TeamID   int    `json:"teamId"`
}

// ========== 공통 DTO ==========

// ErrorResponse 에러 응답
type ErrorResponse struct {
	Message string `json:"message"`
}

// SuccessResponse 성공 응답
type SuccessResponse struct {
	Message string `json:"message"`
}
