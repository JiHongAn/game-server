package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Game 게임 엔티티
type Game struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex"`
	Description string    `json:"description"`
	MaxPlayers  int       `json:"max_players"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate UUID 자동 생성
func (g *Game) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}
