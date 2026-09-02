package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/Newo123/payment-bot-go/internal/domain"
	"github.com/Newo123/payment-bot-go/internal/infrastructure/postgres"
)

type postgresRepository struct {
	pool postgres.Pool
}

func NewPostgresRepository(pool postgres.Pool) *postgresRepository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) FindByTelegramID(ctx context.Context, id int64) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	row := r.pool.QueryRow(ctx, findByIDQuery, id)

	var userModel UserModel
	if err := userModel.Scan(row); err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}

		return domain.User{}, fmt.Errorf("scan error")
	}

	userDomain := modelToDomain(userModel)

	return userDomain, nil
}

func (r *postgresRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	// id, version, phone, language_code, first_name, last_name, username, updated_at, created_at
	row := r.pool.QueryRow(
		ctx,
		insertQuery,
		user.ID,
		user.Version,
		user.Phone,
		user.FirstName,
		user.LastName,
		user.Username,
		user.LanguageCode,
		user.UpdatedAt,
		user.CreatedAt,
	)

	var userModel UserModel
	if err := userModel.Scan(row); err != nil {
		if errors.Is(err, postgres.ErrViolatesUniqueKey) {
			return domain.User{}, domain.ErrAlreadyExists
		}

		return domain.User{}, fmt.Errorf("scan error")
	}

	userDomain := modelToDomain(userModel)

	return userDomain, nil
}
