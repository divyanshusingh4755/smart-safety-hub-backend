package blogs

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/smart-safety-hub/backend/shared"
)

type BlogRepo struct {
	db *sqlx.DB
}

func NewBlogRepo(db *sqlx.DB) *BlogRepo {
	return &BlogRepo{
		db: db,
	}
}

func (r *BlogRepo) SaveBlog(ctx context.Context, request BlogRequestDTO) error {
	query := `
		INSERT INTO blogs (
			title,
			slug,
			type,
			excerpt,
			content,
			featured_image,
			meta_title,
			meta_description,
			is_published,
			published_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, COALESCE($9, FALSE), CASE WHEN COALESCE($9, FALSE) = TRUE THEN CURRENT_TIMESTAMP ELSE NULL END)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		request.Title,
		request.Slug,
		request.Type,
		request.Excerpt,
		request.Content,
		request.FeaturedImage,
		request.MetaTitle,
		request.MetaDescription,
		request.IsPublished,
	)

	if err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *BlogRepo) UpdateBlog(ctx context.Context, blogID string, request UpdateBlogDTO) error {

	var isPublished interface{}

	if request.IsPublished != nil {
		isPublished = *request.IsPublished
	} else {
		isPublished = nil
	}

	query := `
		UPDATE blogs
		SET
			title = COALESCE($1, title),
			slug = COALESCE($2, slug),
			type = COALESCE($3, type),
			excerpt = COALESCE($4, excerpt),
			content = COALESCE($5, content),
			featured_image = COALESCE($6, featured_image),
			meta_title = COALESCE($7, meta_title),
			meta_description = COALESCE($8, meta_description),

			published_at = CASE
				WHEN $9::boolean = TRUE
					AND is_published = FALSE
				THEN CURRENT_TIMESTAMP

				WHEN $9::boolean = FALSE
				THEN NULL

				ELSE published_at
			END,

			is_published = COALESCE(
				$9::boolean,
				is_published
			),

			updated_at = CURRENT_TIMESTAMP

		WHERE id = $10
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		request.Title,
		request.Slug,
		request.Type,
		request.Excerpt,
		request.Content,
		request.FeaturedImage,
		request.MetaTitle,
		request.MetaDescription,
		isPublished,
		blogID,
	)

	if err != nil {
		return shared.PostgresError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("blog not found")
	}

	return nil
}

func (r *BlogRepo) DeleteBlog(ctx context.Context, blogID string) error {
	query := `DELETE FROM blogs WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, blogID)
	if err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *BlogRepo) GetBlogByID(ctx context.Context, blogID string) (*Blog, error) {
	var blog Blog

	query := `
		SELECT
			id,
			title,
			slug,
			type,
			excerpt,
			content,
			featured_image,
			meta_title,
			meta_description,
			is_published,
			published_at,
			created_at,
			updated_at
		FROM blogs
		WHERE id = $1
	`

	if err := r.db.GetContext(ctx, &blog, query, blogID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("blog not found")
		}

		return nil, shared.PostgresError(err)
	}

	return &blog, nil
}

func (r *BlogRepo) getBlogBySlug(ctx context.Context, slug string) (*Blog, error) {
	var blog Blog
	query := `SELECT * FROM blogs WHERE slug = $1 AND is_published = TRUE`

	if err := r.db.GetContext(ctx, &blog, query, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("blog not found")
		}
		return nil, shared.PostgresError(err)
	}
	return &blog, nil
}

func (r *BlogRepo) GetAllBlogs(ctx context.Context) ([]Blog, error) {
	var blogs []Blog
	query := `SELECT * FROM blogs ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &blogs, query); err != nil {
		return nil, shared.PostgresError(err)
	}
	return blogs, nil
}

func (r *BlogRepo) GetPublishedBlogs(ctx context.Context) ([]Blog, error) {
	var blogs []Blog

	query := `
		SELECT
			id,
			title,
			slug,
			type,
			excerpt,
			content,
			featured_image,
			meta_title,
			meta_description,
			is_published,
			published_at,
			created_at,
			updated_at
		FROM blogs
		WHERE is_published = TRUE
		  AND type = 'post'
		ORDER BY published_at DESC
	`

	if err := r.db.SelectContext(ctx, &blogs, query); err != nil {
		return nil, shared.PostgresError(err)
	}

	return blogs, nil
}

func (r *BlogRepo) GetPublishedPages(ctx context.Context) ([]Blog, error) {

	var pages []Blog

	query := `
		SELECT
			id,
			title,
			slug,
			type,
			excerpt,
			content,
			featured_image,
			meta_title,
			meta_description,
			is_published,
			published_at,
			created_at,
			updated_at
		FROM blogs
		WHERE is_published = TRUE
		  AND type = 'page'
		ORDER BY created_at DESC
	`

	if err := r.db.SelectContext(ctx, &pages, query); err != nil {
		return nil, shared.PostgresError(err)
	}

	return pages, nil
}
