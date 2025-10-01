package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntisrangga142/chat/internals/models"
)

type PostRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *PostRepository {
	return &PostRepository{db: db}
}

// Get Following Posts
func (r *PostRepository) GetFollowingPosts(ctx context.Context, followerID int) ([]models.PostFeed, error) {
	query := `
		SELECT 
			p.id, 
			pr.fullname, 
			p.caption, 
			ARRAY_AGG(DISTINCT pi.img) AS images, 
			COUNT(DISTINCT lk.id) AS like_count, 
			COUNT(DISTINCT cm.id) AS comment_count
		FROM posts p
		LEFT JOIN post_imgs pi ON p.id = pi.post_id
		INNER JOIN accounts ac ON p.account_id = ac.id
		INNER JOIN profiles pr ON ac.id = pr.id
		INNER JOIN followers fl ON ac.id = fl.account_id
		LEFT JOIN likes lk ON p.id = lk.post_id AND lk.deleted_at IS NULL
		LEFT JOIN comments cm ON p.id = cm.post_id AND cm.deleted_at IS NULL
		WHERE fl.follower_id = $1 AND fl.deleted_at IS NULL AND p.deleted_at IS NULL
		GROUP BY p.id, pr.fullname, p.caption
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, followerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostFeed
	for rows.Next() {
		var post models.PostFeed
		err := rows.Scan(
			&post.ID,
			&post.Fullname,
			&post.Caption,
			&post.Images,
			&post.LikeCount,
			&post.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Get Post Detail
func (r *PostRepository) GetPostDetail(ctx context.Context, postID int) (*models.PostDetail, error) {
	query := `
		SELECT p.id, p.caption, p.created_at,
		       pr.id, pr.fullname, pr.img,
		       COALESCE(ARRAY_AGG(DISTINCT pi.img) FILTER (WHERE pi.deleted_at IS NULL), '{}') AS images,
		       COUNT(DISTINCT lk.id) FILTER (WHERE lk.deleted_at IS NULL) AS like_count
		FROM posts p
		INNER JOIN profiles pr ON pr.id = p.account_id
		LEFT JOIN post_imgs pi ON p.id = pi.post_id
		LEFT JOIN likes lk ON p.id = lk.post_id
		WHERE p.id = $1 AND p.deleted_at IS NULL
		GROUP BY p.id, pr.id, pr.fullname, pr.img
	`

	var post models.PostDetail
	err := r.db.QueryRow(ctx, query, postID).Scan(
		&post.ID,
		&post.Caption,
		&post.CreatedAt,
		&post.Author.ID,
		&post.Author.Fullname,
		&post.Author.Img,
		&post.Images,
		&post.Likes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed get post detail: %w", err)
	}

	commentQuery := `
		SELECT c.id, c.account_id, pr.fullname, pr.img, c.comment, c.created_at
		FROM comments c
		INNER JOIN profiles pr ON pr.id = c.account_id
		WHERE c.post_id = $1 AND c.deleted_at IS NULL
		ORDER BY c.created_at DESC
		LIMIT 5
	`
	rows, err := r.db.Query(ctx, commentQuery, postID)
	if err != nil {
		return nil, fmt.Errorf("failed get comments: %w", err)
	}
	defer rows.Close()

	var comments []models.CommentPreview
	for rows.Next() {
		var cm models.CommentPreview
		if err := rows.Scan(&cm.ID, &cm.AccountID, &cm.Fullname, &cm.Img, &cm.Comment, &cm.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, cm)
	}
	post.Comments = comments

	return &post, nil
}

// Create Post
func (r *PostRepository) CreatePost(ctx context.Context, req models.CreatePostRequest, accountID int) (*models.Post, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var postID int
	query := `INSERT INTO posts (account_id, caption) VALUES ($1, $2) RETURNING id`
	if err := tx.QueryRow(ctx, query, accountID, req.Caption).Scan(&postID); err != nil {
		return nil, fmt.Errorf("failed to insert post: %w", err)
	}

	// Insert images jika ada
	for _, img := range req.Images {
		_, err := tx.Exec(ctx, `INSERT INTO post_imgs (post_id, img) VALUES ($1, $2)`, postID, img)
		if err != nil {
			return nil, fmt.Errorf("failed to insert post image: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return &models.Post{
		ID:        postID,
		AccountID: accountID,
		Caption:   req.Caption,
		Images:    make([]models.PostImg, 0), // bisa load lagi kalau perlu
	}, nil
}

// Like Post
func (r *PostRepository) CreateLike(ctx context.Context, accountID int, postID int) error {
	query := `
		INSERT INTO likes (account_id, post_id, read, deleted_at)
		VALUES ($1, $2, false, NULL)
		ON CONFLICT (account_id, post_id)
		DO UPDATE SET deleted_at = NULL, read = false
	`
	_, err := r.db.Exec(ctx, query, accountID, postID)
	if err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}
	return nil
}

// Unlike Post
func (r *PostRepository) DeleteLike(ctx context.Context, accountID, postID int) error {
	query := `
		UPDATE likes SET deleted_at = NOW()
		WHERE account_id=$1 AND post_id=$2 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, accountID, postID)
	if err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}
	return nil
}

// Create Comment Post
func (r *PostRepository) CreateComment(ctx context.Context, accountID int, req models.CreateCommentRequest) error {
	query := `
		INSERT INTO comments (account_id, post_id, comment, read)
		VALUES ($1, $2, $3, false)
	`
	_, err := r.db.Exec(ctx, query, accountID, req.PostID, req.Comment)
	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}
	return nil
}

// Get Comment Post
func (r *PostRepository) GetAllCommentsByPost(ctx context.Context, postID int) ([]models.Comment, error) {
	query := `
		SELECT id, account_id, post_id, comment
		FROM comments
		WHERE post_id=$1 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.AccountID, &c.PostID, &c.Comment); err == nil {
			comments = append(comments, c)
		}
	}

	return comments, nil
}
