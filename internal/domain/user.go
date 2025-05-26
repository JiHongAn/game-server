package domain

import "time"

// User 사용자 엔티티
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Game 게임 엔티티
type Game struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MaxPlayers  int       `json:"max_players"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Match 매치 엔티티
type Match struct {
	ID        string    `json:"id"`
	GameID    string    `json:"game_id"`
	Players   []string  `json:"players"`
	Status    string    `json:"status"` // waiting, in_progress, completed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
