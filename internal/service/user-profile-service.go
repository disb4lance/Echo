package service

import (
	"echo/internal/domain/entity"
	"echo/internal/service/dto"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserProfileRepository interface {
	Create(profile *entity.UserProfile) error
	GetById(userId uuid.UUID) (*entity.UserProfile, error)
	Update(profile *entity.UserProfile) error
}

var (
	ErrProfileNotFound      = errors.New("profile not found")
	ErrProfileAlreadyExists = errors.New("profile already exists")
	ErrInvalidBirthDate     = errors.New("invalid birth date")
	ErrInvalidGender        = errors.New("invalid gender")
	ErrUnauthorized         = errors.New("unauthorized")
)

type UserProfileService struct {
	repo UserProfileRepository
}

func NewUserProfileService(repo UserProfileRepository) *UserProfileService {
	return &UserProfileService{
		repo: repo,
	}
}

func (s *UserProfileService) CreateProfile(request *dto.UserProfileRequest) (*dto.UserProfileResponse, error) {

	userID := uuid.New() //TODO брать из токена

	// existing, err := s.repo.GetById(userID)
	// if err != nil {
	// 	return nil, err
	// }
	// if existing != nil {
	// 	return nil, ErrProfileAlreadyExists
	// }

	profile := dto.ToEntity(request, userID)

	if err := s.repo.Create(profile); err != nil {
		return nil, err
	}
	return dto.ToResponse(profile), nil
}

func (s *UserProfileService) GetProfile(userID uuid.UUID) (*dto.UserProfileResponse, error) {
	if userID == uuid.Nil {
		return nil, ErrUnauthorized
	}

	profile, err := s.repo.GetById(userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, ErrProfileNotFound
	}

	return dto.ToResponse(profile), nil
}

func (s *UserProfileService) UpdateProfile(userID uuid.UUID, request *dto.UserProfileRequest) (*dto.UserProfileResponse, error) {
	if userID == uuid.Nil {
		return nil, ErrUnauthorized
	}

	profile, err := s.repo.GetById(userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, ErrProfileNotFound
	}

	if request.BirthDate.After(time.Now()) {
		return nil, ErrInvalidBirthDate
	}
	if !request.Gender.IsValid() {
		return nil, ErrInvalidGender
	}

	profile.Name = request.Name
	profile.BirthDate = request.BirthDate
	profile.Gender = request.Gender
	profile.Description = request.Description
	profile.UpdatedAt = time.Now()

	if err := s.repo.Update(profile); err != nil {
		return nil, err
	}

	return dto.ToResponse(profile), nil
}
