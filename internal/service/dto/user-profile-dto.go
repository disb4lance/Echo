package dto

import (
	"echo/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type Gender = entity.Gender
type City = entity.City

const (
	PreferNotToSay Gender = iota
	Male
	Female
	Other
)

const (
	Moscow City = iota
	Saint_Petersburg
)

type UserProfileRequest struct {
	Name        string    `json:"name" validate:"required,min=2,max=100"`
	BirthDate   time.Time `json:"birth_date" validate:"required"`
	Gender      Gender    `json:"gender" validate:"required"`
	Description string    `json:"description" validate:"max=500"`
}

type UserProfileResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Age         int       `json:"age"`
	Gender      Gender    `json:"gender"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Конвертеры
func ToResponse(profile *entity.UserProfile) *UserProfileResponse {
	if profile == nil {
		return nil
	}

	return &UserProfileResponse{
		UserID:      profile.UserID,
		Name:        profile.Name,
		Age:         profile.CalculateAge(),
		Gender:      profile.Gender,
		Description: profile.Description,
		IsActive:    profile.IsActive,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}
}

func ToEntity(request *UserProfileRequest, userID uuid.UUID) *entity.UserProfile {
	now := time.Now()
	return &entity.UserProfile{
		UserID:      userID,
		Name:        request.Name,
		BirthDate:   request.BirthDate,
		Gender:      request.Gender,
		Description: request.Description,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
