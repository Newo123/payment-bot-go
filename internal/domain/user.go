package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type User struct {
	ID           int64
	Version      int64
	Phone        string
	LanguageCode string
	FirstName    string
	LastName     string
	Username     string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

func NewUser(
	id int64,
	version int64,
	phone string,
	languageCode string,
	firstName string,
	lastName string,
	username string,
	updatedAt time.Time,
	createdAt time.Time,
) User {
	return User{
		ID:           id,
		Version:      version,
		Phone:        phone,
		LanguageCode: languageCode,
		FirstName:    firstName,
		LastName:     lastName,
		Username:     username,
		UpdatedAt:    updatedAt,
		CreatedAt:    createdAt,
	}
}

func CreateUser(
	id int64,
	phone string,
	languageCode string,
	firstName string,
	lastName string,
	username string,
) (User, error) {
	now := time.Now()

	user := NewUser(
		id,
		1,
		phone,
		languageCode,
		firstName,
		lastName,
		username,
		now,
		now,
	)

	if err := user.validate(); err != nil {
		return User{}, fmt.Errorf("validation err: %w", err)
	}

	return user, nil
}

func (u User) validate() error {
	// Валидация Phone (если заполнен)
	if u.Phone != "" {
		phoneRegex := regexp.MustCompile(`^(\+7|8)\d{10}$`)
		if !phoneRegex.MatchString(u.Phone) {
			return fmt.Errorf("phone must be in format +7XXXXXXXXXX or 8XXXXXXXXXX: %w", ErrInvalidArgument)
		}
	}

	// Валидация TelegramID (если заполнен)
	if u.ID != 0 && u.ID <= 0 {
		return fmt.Errorf("telegram_id must be positive: %w", ErrInvalidArgument)
	}

	// Валидация LanguageCode (если заполнен)
	if u.LanguageCode != "" {
		if len(u.LanguageCode) != 2 {
			return fmt.Errorf("language_code must be 2 characters (ISO 639-1): %w", ErrInvalidArgument)
		}
	}

	// Валидация FirstName (если заполнен)
	if u.FirstName != "" {
		trimmed := strings.TrimSpace(u.FirstName)
		if trimmed == "" {
			return fmt.Errorf("first_name cannot be empty or only spaces: %w", ErrInvalidArgument)
		}
		if len(trimmed) > 64 {
			return fmt.Errorf("first_name must be less than 64 characters: %w", ErrInvalidArgument)
		}
	}

	// Валидация LastName (если заполнен)
	if u.LastName != "" {
		trimmed := strings.TrimSpace(u.LastName)
		if trimmed != "" && len(trimmed) > 64 {
			return fmt.Errorf("last_name must be less than 64 characters: %w", ErrInvalidArgument)
		}
	}

	// Валидация Username (если заполнен)
	if u.Username != "" {
		trimmed := strings.TrimSpace(u.Username)
		if len(trimmed) < 3 {
			return fmt.Errorf("username must be at least 3 characters: %w", ErrInvalidArgument)
		}
		if len(trimmed) > 32 {
			return fmt.Errorf("username must be less than 32 characters: %w", ErrInvalidArgument)
		}
		usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
		if !usernameRegex.MatchString(trimmed) {
			return fmt.Errorf("username can only contain letters, numbers and underscore: %w", ErrInvalidArgument)
		}
	}

	return nil
}
