package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntisrangga142/chat/internals/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Get Profile
func (ur *UserRepository) GetProfile(ctx context.Context, uid int) (*models.Profile, error) {
	sql := `
		SELECT p.fullname, p.phone, p.img
		FROM profiles p
		WHERE p.id = $1
	`

	row := ur.db.QueryRow(ctx, sql, uid)

	var profile models.Profile
	err := row.Scan(
		&profile.FullName,
		&profile.PhoneNumber,
		&profile.Img,
	)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

// Update Profile
func (ur *UserRepository) UpdateProfile(ctx context.Context, uid int, updates map[string]any) error {
	if len(updates) == 0 {
		return nil
	}

	setClauses := []string{}
	args := []any{}
	i := 1

	for col, val := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}

	// tambah updated_at
	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf(`
		UPDATE profiles
		SET %s
		WHERE id = $%d
	`, strings.Join(setClauses, ", "), i)

	args = append(args, uid)

	_, err := ur.db.Exec(ctx, query, args...)
	return err
}

// Follow user
func (r *UserRepository) Follow(ctx context.Context, accountID, followerID int) error {
	query := `
		INSERT INTO followers (account_id, follower_id, read)
		VALUES ($1, $2, false)
		ON CONFLICT (account_id, follower_id) 
		DO UPDATE SET deleted_at = NULL
	`
	_, err := r.db.Exec(ctx, query, accountID, followerID)
	if err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}
	return nil
}

// Unfollow user (soft delete)
func (r *UserRepository) Unfollow(ctx context.Context, accountID, followerID int) error {
	query := `
		UPDATE followers
		SET deleted_at = NOW()
		WHERE account_id=$1 AND follower_id=$2
	`
	_, err := r.db.Exec(ctx, query, accountID, followerID)
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}
	return nil
}

// Get Followers
func (r *UserRepository) GetFollowers(ctx context.Context, accountID int) ([]models.Follow, error) {
	query := `
		SELECT p.id, p.fullname, p.img
		FROM followers f
		JOIN profiles p ON p.id = f.follower_id
		WHERE f.account_id = $1 AND f.deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []models.Follow
	for rows.Next() {
		var p models.Follow
		if err := rows.Scan(&p.ID, &p.Fullname, &p.Img); err == nil {
			profiles = append(profiles, p)
		}
	}

	return profiles, nil
}

// Get Following
func (r *UserRepository) GetFollowing(ctx context.Context, followerID int) ([]models.Follow, error) {
	query := `
		SELECT p.id, p.fullname, p.img
		FROM followers f
		JOIN profiles p ON p.id = f.account_id
		WHERE f.follower_id = $1 AND f.deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query, followerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []models.Follow
	for rows.Next() {
		var p models.Follow
		if err := rows.Scan(&p.ID, &p.Fullname, &p.Img); err == nil {
			profiles = append(profiles, p)
		}
	}

	return profiles, nil
}
