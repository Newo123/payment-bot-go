package users

import (
	"context"
	"fmt"

	"github.com/Newo123/payment-bot-go/internal/domain"
)

type UseCase interface {
	FindByID(
		ctx context.Context,
		telegramID int64,
	) (domain.User, error)
	Create(
		ctx context.Context,
		params CreateParams,
	) (domain.User, error)
	Update(
		ctx context.Context,
		telegramID int64,
		params UpdateParams,
	) (domain.User, error)
}

type UpdateParams struct {
	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
}

type CreateParams struct {
	ID           int64
	Phone        string
	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
}

// type UpdateParams struct {
// 	ID uuid.UUID
// 	Version int64

// }

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *useCase {
	return &useCase{
		repo: repo,
	}
}

func (s *useCase) FindByID(ctx context.Context, id int64) (domain.User, error) {
	user, err := s.repo.FindByTelegramID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s *useCase) Create(ctx context.Context, params CreateParams) (domain.User, error) {
	userDomain, err := domain.CreateUser(
		params.ID,
		params.Phone,
		params.LanguageCode,
		params.FirstName,
		params.LastName,
		params.Username,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user domain: %w", err)
	}

	user, err := s.repo.Create(ctx, userDomain)
	if err != nil {
		return domain.User{}, fmt.Errorf("save user in repo: %w", err)
	}

	return user, nil
}

func (s *useCase) Update(ctx context.Context, id int64, params UpdateParams) (domain.User, error) {

	return domain.User{}, nil
}
