package service

import (
	"game-server/pkg/errors"
)

type GameService struct{}

func NewGameService() *GameService {
	return &GameService{}
}

func (service *GameService) GetGames() (bool, *errors.AppError) {
	return true, nil
}
