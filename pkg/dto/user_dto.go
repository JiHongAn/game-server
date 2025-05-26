package dto

type UserResponse struct {
	ID         string  `json:"id"`
	Nickname   string  `json:"nickname"`
	ProfileKey *string `json:"profileKey"`
}

type UserUpdateRequest struct {
	Nickname   string `json:"nickname"`
	ProfileKey string `json:"profileKey"`
}
