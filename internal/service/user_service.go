package service

import (
	"game-server/internal/domain"
	"game-server/internal/pkg/errors"
)

// UserService 사용자 비즈니스 로직
type UserService struct {
	// ORM을 직접 사용하거나 DB 연결을 여기서 처리
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUserByID(id string) (*domain.User, *errors.AppError) {
	// TODO: ORM을 사용하여 사용자 조회
	// 예: db.First(&user, "id = ?", id)
	return nil, errors.NotFound()
}

func (s *UserService) CreateUser(user *domain.User) *errors.AppError {
	// 비즈니스 로직 검증
	if user.Username == "" {
		return errors.BadRequestWithMessage("사용자명은 필수입니다")
	}
	if user.Email == "" {
		return errors.BadRequestWithMessage("이메일은 필수입니다")
	}

	// TODO: ORM을 사용하여 이메일 중복 확인
	// 예: var existingUser domain.User
	//     if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
	//         return errors.BadRequestWithMessage("이미 존재하는 이메일입니다")
	//     }

	// TODO: ORM을 사용하여 사용자 생성
	// 예: db.Create(user)
	return nil
}

// GameService 게임 비즈니스 로직
type GameService struct {
	// ORM을 직접 사용하거나 DB 연결을 여기서 처리
}

func NewGameService() *GameService {
	return &GameService{}
}

func (s *GameService) GetGames() ([]*domain.Game, *errors.AppError) {
	// TODO: ORM을 사용하여 게임 목록 조회
	// 예: var games []*domain.Game
	//     db.Find(&games)
	//     return games, nil

	// 임시로 빈 배열 반환
	return []*domain.Game{}, nil
}
