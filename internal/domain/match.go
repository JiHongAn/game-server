package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
