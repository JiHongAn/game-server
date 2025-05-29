package service

import (
	"game-server/internal/dto"
	"time"
)

type GameService struct {
}

type Game struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"`
	PlayerCount int                    `json:"playerCount"`
	Teams       []dto.Team             `json:"teams"`
	CreatedAt   time.Time              `json:"createdAt"`
	StartedAt   *time.Time             `json:"startedAt,omitempty"`
	EndedAt     *time.Time             `json:"endedAt,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

const (
	GAME_STATUS_WAITING  = "waiting"
	GAME_STATUS_STARTING = "starting"
	GAME_STATUS_PLAYING  = "playing"
	GAME_STATUS_ENDED    = "ended"
)

func NewGameService() *GameService {
	return &GameService{}
}
