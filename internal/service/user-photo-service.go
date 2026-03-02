package service

import (
	"context"
	"echo/internal/adapter/postgres"
	"echo/internal/domain/entity"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type MinioStorage interface {
	AddFiles(ctx context.Context, userID string, files []*multipart.FileHeader) ([]string, error)
	DeleteUserFiles(ctx context.Context, userID string) error
	ReplaceUserFiles(ctx context.Context, userID string, files []*multipart.FileHeader) ([]string, error)
	GetUserFiles(ctx context.Context, userID string) ([]string, error)
}

type UserPhotoService struct {
	txManager  *postgres.TxManager
	minioStore MinioStorage
}

func NewUserPhotoService(txManager *postgres.TxManager, minioStore MinioStorage) *UserPhotoService {
	return &UserPhotoService{
		txManager:  txManager,
		minioStore: minioStore,
	}
}

func (s *UserPhotoService) AddUserPhotos(ctx context.Context, userID uuid.UUID, files []*multipart.FileHeader) ([]entity.UserPhoto, error) {
	urls, err := s.minioStore.AddFiles(ctx, userID.String(), files)
	if err != nil {
		return nil, err
	}

	photos := make([]entity.UserPhoto, 0, len(urls))
	for i, url := range urls {
		photos = append(photos, entity.UserPhoto{
			ID:        uuid.New(),
			UserID:    userID,
			URL:       url,
			Position:  i,
			CreatedAt: time.Now(),
		})
	}

	err = s.txManager.WithTx(ctx, func(tx postgres.DBTX) error {
		repo := postgres.NewUserPhotoRepo(tx)
		return repo.CreateMany(photos)
	})
	if err != nil {
		return nil, err
	}

	return photos, nil
}

func (s *UserPhotoService) GetUserPhotos(ctx context.Context, userID uuid.UUID) ([]entity.UserPhoto, error) {
	repo := postgres.NewUserPhotoRepo(s.txManager.Pool())
	return repo.GetByUserID(userID)
}

func (s *UserPhotoService) ReplaceUserPhotos(ctx context.Context, userID uuid.UUID, files []*multipart.FileHeader) ([]entity.UserPhoto, error) {
	urls, err := s.minioStore.ReplaceUserFiles(ctx, userID.String(), files)
	if err != nil {
		return nil, err
	}

	photos := make([]entity.UserPhoto, 0, len(urls))
	for i, url := range urls {
		photos = append(photos, entity.UserPhoto{
			ID:        uuid.New(),
			UserID:    userID,
			URL:       url,
			Position:  i,
			CreatedAt: time.Now(),
		})
	}

	err = s.txManager.WithTx(ctx, func(uow *postgres.UnitOfWork) error {
		if err := uow.UserPhotoRepo.DeleteByUserID(userID); err != nil {
			return err
		}
		return uow.UserPhotoRepo.CreateMany(photos)
	})
	if err != nil {
		return nil, err
	}

	return photos, nil
}

func (s *UserPhotoService) DeleteUserPhotos(ctx context.Context, userID uuid.UUID) error {

	if err := s.minioStore.DeleteUserFiles(ctx, userID.String()); err != nil {
		return err
	}

	return s.txManager.WithTx(ctx, func(uow *postgres.UnitOfWork) error {

		if err := uow.UserPhotoRepo.DeleteByUserID(userID); err != nil {
			return err
		}

		return nil
	})
}
