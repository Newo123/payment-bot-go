package users

import (
	"context"
	"time"

	"github.com/Newo123/payment-bot-go/internal/domain"
	"github.com/Newo123/payment-bot-go/internal/infrastructure/postgres"
)

type Repository interface {
	FindByTelegramID(
		ctx context.Context,
		telegramID int64,
	) (domain.User, error)
	Create(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

var (
	findByIDQuery = `
	SELECT id, version, phone, first_name, last_name, username, language_code, updated_at, created_at
	FROM bot.users
	WHERE id=$1;
	`
	insertQuery = `
	INSERT INTO bot.users (id, version, phone, first_name, last_name, username, language_code, updated_at, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, version, phone, first_name, last_name, username, language_code, updated_at, created_at;
	`
)

type UserModel struct {
	ID           int64
	Version      int64
	Phone        string
	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

// Scan заполняет поля модели из результата запроса к БД.
func (m *UserModel) Scan(row postgres.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.Phone,
		&m.FirstName,
		&m.LastName,
		&m.Username,
		&m.LanguageCode,
		&m.UpdatedAt,
		&m.CreatedAt,
	)
}

// modelToDomain конвертирует модель БД в доменный объект.
func modelToDomain(model UserModel) domain.User {
	return domain.NewUser(
		model.ID,
		model.Version,
		model.Phone,
		model.FirstName,
		model.LastName,
		model.Username,
		model.LanguageCode,
		model.UpdatedAt,
		model.CreatedAt,
	)
}

// modelsToDomains конвертирует список моделей БД в список доменных объектов.
func modelsToDomains(models []UserModel) []domain.User {
	domains := make([]domain.User, len(models))

	for i, model := range models {
		domains[i] = modelToDomain(model)
	}

	return domains
}
