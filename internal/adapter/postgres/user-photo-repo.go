package postgres

import (
	"context"
	"echo/internal/domain/entity"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)

	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type UserPhotoRepo struct {
	db DBTX
}

func NewUserPhotoRepo(db DBTX) *UserPhotoRepo {
	return &UserPhotoRepo{
		db: db,
	}
}

func (r *UserPhotoRepo) CreateMany(photos []entity.UserPhoto) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, photo := range photos {

		if photo.ID == uuid.Nil {
			photo.ID = uuid.New()
		}

		if photo.CreatedAt.IsZero() {
			photo.CreatedAt = time.Now()
		}

		_, err := r.db.Exec(ctx,
			`INSERT INTO user_photos (id, user_id, url, position, created_at)
			 VALUES ($1, $2, $3, $4, $5)`,
			photo.ID,
			photo.UserID,
			photo.URL,
			photo.Position,
			photo.CreatedAt,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserPhotoRepo) GetByUserID(userID uuid.UUID) ([]entity.UserPhoto, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, url, position, created_at
		 FROM user_photos
		 WHERE user_id = $1
		 ORDER BY position ASC`,
		userID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []entity.UserPhoto

	for rows.Next() {

		var photo entity.UserPhoto

		err := rows.Scan(
			&photo.ID,
			&photo.UserID,
			&photo.URL,
			&photo.Position,
			&photo.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		photos = append(photos, photo)
	}

	return photos, nil
}

func (r *UserPhotoRepo) DeleteByUserID(userID uuid.UUID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		`DELETE FROM user_photos
		 WHERE user_id = $1`,
		userID,
	)

	return err
}
