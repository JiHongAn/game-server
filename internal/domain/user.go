package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 사용자 엔티티
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate UUID 자동 생성
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

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

// Match 매치 엔티티
type Match struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	GameID    string    `json:"game_id"`
	Players   string    `json:"players" gorm:"type:text"` // JSON 문자열로 저장
	Status    string    `json:"status"`                   // waiting, in_progress, completed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate UUID 자동 생성
func (m *Match) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
