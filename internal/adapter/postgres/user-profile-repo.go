package postgres

import (
	"context"
	"echo/internal/domain/entity"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserProfileRepo struct {
	db *pgxpool.Pool
}

func NewUserProfileRepo(db *pgxpool.Pool) *UserProfileRepo {
	return &UserProfileRepo{db: db}
}

func (r *UserProfileRepo) Create(profile *entity.UserProfile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		`INSERT INTO user_profiles (user_id, name, birth_date, gender, description, is_active, created_at, updated_at)
    VALUES ($1, $2, $3, 'male', $4, $5, $6, $7)`,
		profile.UserID, profile.Name, profile.BirthDate, profile.Description,
		profile.IsActive, profile.CreatedAt, profile.UpdatedAt,
	)
	return err
}

func (r *UserProfileRepo) GetById(userId uuid.UUID) (*entity.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var profile entity.UserProfile
	row := r.db.QueryRow(ctx,
		`SELECT user_id, name, birth_date, gender, description, is_active, created_at, updated_at
		 FROM user_profiles 
		 WHERE user_id = $1`,
		userId,
	)

	err := row.Scan(&profile.UserID, &profile.Name, &profile.BirthDate, &profile.Gender, &profile.Description,
		&profile.IsActive, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

func (r *UserProfileRepo) Update(profile *entity.UserProfile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		`UPDATE user_profiles 
         SET name = $1, birth_date = $2, gender = $3, description = $4, 
             is_active = $5, updated_at = $6
         WHERE user_id = $7`,
		profile.Name, profile.BirthDate, profile.Gender, profile.Description,
		profile.IsActive, profile.UpdatedAt, profile.UserID,
	)
	return err
}
