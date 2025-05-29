package service

import (
	"encoding/json"
	"fmt"
	"game-server/internal/dto"
	"game-server/internal/pkg/database"
	"time"

	"github.com/google/uuid"
)

type MatchService struct{}

func NewMatchService() *MatchService {
	return &MatchService{}
}

const (
	INVITE_EXPIRE_MINUTES = 5
)

// CreateMatch 새로운 매치 생성
func (s *MatchService) CreateMatch(hostID, gameID string, maxPlayers int) (*dto.MatchInfo, error) {
	if gameID == "" || maxPlayers < 2 {
		return nil, fmt.Errorf("gameId and valid maxPlayers are required")
	}

	matchID := uuid.New().String()

	// 매치 정보 생성
	matchInfo := &dto.MatchInfo{
		MatchID:    matchID,
		GameID:     gameID,
		HostID:     hostID,
		Status:     "waiting",
		Players:    []dto.MatchPlayer{},
		MaxPlayers: maxPlayers,
		CreatedAt:  time.Now(),
	}

	// Redis에 매치 정보 저장
	matchJSON, _ := json.Marshal(matchInfo)
	if err := database.HSet("matches", matchID, string(matchJSON)); err != nil {
		return nil, fmt.Errorf("failed to create match: %w", err)
	}

	// 사용자를 매치에 연결
	if err := database.HSet("user:matches", hostID, matchID); err != nil {
		return nil, fmt.Errorf("failed to link user to match: %w", err)
	}

	return matchInfo, nil
}

// InviteFriends 친구들을 매치에 초대
func (s *MatchService) InviteFriends(hostID, matchID string, friendIds []string) (*dto.InviteFriendResponse, error) {
	if matchID == "" || len(friendIds) == 0 {
		return nil, fmt.Errorf("matchId and friendIds are required")
	}

	// 매치 정보 확인
	matchInfo, err := s.GetMatchInfo(matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found")
	}

	// 호스트 권한 확인
	if matchInfo.HostID != hostID {
		return nil, fmt.Errorf("only host can invite friends")
	}

	var invitedIds []string
	var failedIds []string
	expiresAt := time.Now().Add(INVITE_EXPIRE_MINUTES * time.Minute).Unix()

	for _, friendID := range friendIds {
		// 친구가 온라인인지 확인
		_, err := database.HGet("user:sockets", friendID)
		if err != nil {
			failedIds = append(failedIds, friendID)
			continue
		}

		// 초대 정보 저장 (TTL 설정)
		inviteKey := fmt.Sprintf("invite:%s:%s", matchID, friendID)
		inviteData := dto.MatchInvitation{
			MatchID:   matchID,
			GameID:    matchInfo.GameID,
			HostID:    hostID,
			ExpiresAt: expiresAt,
			Message:   fmt.Sprintf("You are invited to join match for %s", matchInfo.GameID),
		}
		inviteJSON, _ := json.Marshal(inviteData)

		if err := database.Set(inviteKey, string(inviteJSON), INVITE_EXPIRE_MINUTES*time.Minute); err != nil {
			failedIds = append(failedIds, friendID)
			continue
		}

		invitedIds = append(invitedIds, friendID)
	}

	return &dto.InviteFriendResponse{
		MatchID:    matchID,
		InvitedIds: invitedIds,
		FailedIds:  failedIds,
		Message:    fmt.Sprintf("Invited %d friends, %d failed", len(invitedIds), len(failedIds)),
	}, nil
}

// RespondInvite 초대에 응답
func (s *MatchService) RespondInvite(userID, matchID, response string) (*dto.InviteResponseResponse, error) {
	if matchID == "" || (response != "accept" && response != "decline") {
		return nil, fmt.Errorf("matchId and valid response (accept/decline) are required")
	}

	// 초대 정보 확인
	inviteKey := fmt.Sprintf("invite:%s:%s", matchID, userID)
	inviteData, err := database.Get(inviteKey)
	if err != nil || inviteData == "" {
		return nil, fmt.Errorf("invitation not found or expired")
	}

	// 초대 삭제
	database.Del(inviteKey)

	result := &dto.InviteResponseResponse{
		MatchID:  matchID,
		UserID:   userID,
		Response: response,
	}

	if response == "decline" {
		result.Message = "Invitation declined"
		return result, nil
	}

	// accept인 경우 매치에 참가
	matchInfo, err := s.GetMatchInfo(matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found")
	}

	// 매치가 가득 찼는지 확인
	if len(matchInfo.Players) >= matchInfo.MaxPlayers {
		return nil, fmt.Errorf("match is full")
	}

	// 플레이어 추가
	now := time.Now()
	newPlayer := dto.MatchPlayer{
		UserID:   userID,
		Status:   "joined",
		JoinedAt: &now,
	}
	matchInfo.Players = append(matchInfo.Players, newPlayer)

	// 매치 정보 업데이트
	if err := s.UpdateMatchInfo(matchInfo); err != nil {
		return nil, fmt.Errorf("failed to update match")
	}

	// 사용자를 매치에 연결
	database.HSet("user:matches", userID, matchID)

	result.Message = "Successfully joined the match"
	return result, nil
}

// StartMatch 매치 시작
func (s *MatchService) StartMatch(hostID, matchID string) (*dto.StartMatchResponse, error) {
	// 매치 정보 확인
	matchInfo, err := s.GetMatchInfo(matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found")
	}

	// 호스트 권한 확인
	if matchInfo.HostID != hostID {
		return nil, fmt.Errorf("only host can start the match")
	}

	// 최소 플레이어 수 확인
	if len(matchInfo.Players) < 2 {
		return nil, fmt.Errorf("need at least 2 players to start")
	}

	// 매치 상태 업데이트
	now := time.Now()
	matchInfo.Status = "starting"
	matchInfo.StartedAt = &now

	// 팀 생성 (간단하게 절반씩 나누기)
	teams := s.CreateTeams(matchInfo.Players)

	// 매치 정보 업데이트
	if err := s.UpdateMatchInfo(matchInfo); err != nil {
		return nil, fmt.Errorf("failed to update match")
	}

	return &dto.StartMatchResponse{
		MatchID: matchID,
		GameID:  matchInfo.GameID,
		Teams:   teams,
		Message: "Match started!",
	}, nil
}

// LeaveMatch 매치 나가기
func (s *MatchService) LeaveMatch(userID string) error {
	// 사용자의 현재 매치 확인
	matchID, err := database.HGet("user:matches", userID)
	if err != nil || matchID == "" {
		return fmt.Errorf("you are not in any match")
	}

	return s.RemovePlayerFromMatch(userID, matchID)
}

// GetMatchInfo 매치 정보 조회
func (s *MatchService) GetMatchInfo(matchID string) (*dto.MatchInfo, error) {
	matchData, err := database.HGet("matches", matchID)
	if err != nil || matchData == "" {
		return nil, fmt.Errorf("match not found")
	}

	var matchInfo dto.MatchInfo
	if err := json.Unmarshal([]byte(matchData), &matchInfo); err != nil {
		return nil, fmt.Errorf("invalid match data")
	}

	return &matchInfo, nil
}

// UpdateMatchInfo 매치 정보 업데이트
func (s *MatchService) UpdateMatchInfo(matchInfo *dto.MatchInfo) error {
	matchJSON, err := json.Marshal(matchInfo)
	if err != nil {
		return err
	}
	return database.HSet("matches", matchInfo.MatchID, string(matchJSON))
}

// RemovePlayerFromMatch 매치에서 플레이어 제거
func (s *MatchService) RemovePlayerFromMatch(userID, matchID string) error {
	// 사용자-매치 연결 제거
	database.HDel("user:matches", userID)

	// 매치에서 플레이어 제거
	matchInfo, err := s.GetMatchInfo(matchID)
	if err != nil {
		return err
	}

	// 플레이어 목록에서 제거
	var newPlayers []dto.MatchPlayer
	for _, player := range matchInfo.Players {
		if player.UserID != userID {
			newPlayers = append(newPlayers, player)
		}
	}
	matchInfo.Players = newPlayers

	// 매치 정보 업데이트 또는 삭제
	if len(newPlayers) == 0 {
		database.HDel("matches", matchID)
	} else {
		s.UpdateMatchInfo(matchInfo)
	}

	return nil
}

// CreateTeams 팀 생성
func (s *MatchService) CreateTeams(players []dto.MatchPlayer) []dto.Team {
	teams := make([]dto.Team, 2)
	teams[0] = dto.Team{ID: 1, Players: []string{}}
	teams[1] = dto.Team{ID: 2, Players: []string{}}

	for i, player := range players {
		teamIndex := i % 2
		teams[teamIndex].Players = append(teams[teamIndex].Players, player.UserID)
	}

	return teams
}

// GetMatchPlayers 매치의 플레이어 목록 조회
func (s *MatchService) GetMatchPlayers(matchID string) ([]dto.MatchPlayer, error) {
	matchInfo, err := s.GetMatchInfo(matchID)
	if err != nil {
		return nil, err
	}
	return matchInfo.Players, nil
}
