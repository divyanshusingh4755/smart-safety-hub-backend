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

func upsertPageConfigTx(ctx context.Context, tx *sqlx.Tx, blogID string, request PageConfigRequestDTO) error {
	query := `
		INSERT INTO blog_page_configs (
			blog_id, layout, market, hero, trust, components, sections, final_cta, products	
		)
		SELECT id, $2, $3, $4, $5, $6, $7, $8, $9
		FROM blogs WHERE id  = $1 AND type = 'page' ON CONFLICT (blog_id) Do UPDATE SET
			layout = EXCLUDED.layout,
            market = EXCLUDED.market,
            hero = EXCLUDED.hero,
            trust = EXCLUDED.trust,
            components = EXCLUDED.components,
            sections = EXCLUDED.sections,
            final_cta = EXCLUDED.final_cta,
            products = EXCLUDED.products,
            updated_at = CURRENT_TIMESTAMP
		`

	result, err := tx.ExecContext(
		ctx,
		query,
		blogID,
		request.Layout,
		string(request.Market),
		string(request.Hero),
		string(request.Trust),
		string(request.Components),
		string(request.Sections),
		string(request.FinalCTA),
		string(request.Products),
	)
	if err != nil {
		return shared.PostgresError(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
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

func (r *BlogRepo) GetBlogBySlug(ctx context.Context, slug string) (*Blog, error) {
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

func (r *BlogRepo) UpsertPageConfig(ctx context.Context, blogID string, request PageConfigRequestDTO) (err error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := upsertPageConfigTx(ctx, tx, blogID, request); err != nil {
		return err
	}

	// Replace component insertions
	if _, err := tx.ExecContext(ctx, `DELETE FROM blog_page_insertions WHERE blog_id = $1`, blogID); err != nil {
		return err
	}

	for i, item := range request.Insertions {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO blog_page_insertions(blog_id, component, target, position, sort_order)
			VALUES ($1, $2, $3, $4, $5)
		`, blogID, item.Component, item.Target, item.Position, i)

		if err != nil {
			return shared.PostgresError(err)
		}
	}

	// Replace internal links
	if _, err := tx.ExecContext(ctx, `DELETE FROM blog_page_internal_links WHERE blog_id = $1`, blogID); err != nil {
		return err
	}

	for i, item := range request.InternalLinks {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO blog_page_internal_links (blog_id, text, href, max_occurrences, sort_order)
			VALUES ($1, $2, $3, $4, $5)
		`, blogID, item.Text, item.Href, item.MaxOccurrences, i)

		if err != nil {
			return shared.PostgresError(err)
		}
	}

	// Replace content images.
	if _, err := tx.ExecContext(ctx, `DELETE FROM blog_page_images WHERE blog_id = $1`, blogID); err != nil {
		return err
	}

	for i, item := range request.Images {
		_, err := tx.ExecContext(ctx, `
            INSERT INTO blog_page_images (blog_id, target, position, src, alt, sort_order)
            VALUES ($1, $2, $3, $4, $5, $6)
        `, blogID, item.Target, item.Position, item.Src, item.Alt, i)

		if err != nil {
			return shared.PostgresError(err)
		}
	}

	return tx.Commit()
}

func (r *BlogRepo) GetPageConfig(ctx context.Context, blogID string) (*BlogPageConfig, error) {
	var config BlogPageConfig
	query := ` SELECT id, blog_id, layout, market, hero, trust, components, sections, final_cta, products, created_at, updated_at FROM blog_page_configs WHERE blog_id = $1 `

	err := r.db.GetContext(ctx, &config, query, blogID)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (r *BlogRepo) GetPageInsertions(ctx context.Context, blogID string) ([]BlogPageInsertion, error) {
	items := make([]BlogPageInsertion, 0)
	err := r.db.SelectContext(ctx, &items, `
		SELECT component, target, position, sort_order FROM blog_page_insertions WHERE blog_id = $1 ORDER BY sort_order ASC, id ASC
	`, blogID)

	return items, err
}

func (r *BlogRepo) GetPageInternalLinks(ctx context.Context, blogID string) ([]BlogPageInternalLink, error) {
	items := make([]BlogPageInternalLink, 0)
	err := r.db.SelectContext(ctx, &items, `
		SELECT text, href, max_occurrences, sort_order FROM blog_page_internal_links WHERE blog_id = $1 ORDER BY sort_order ASC, id ASC
	`, blogID)

	return items, err
}

func (r *BlogRepo) GetPageImages(ctx context.Context, blogID string) ([]BlogPageImage, error) {
	items := make([]BlogPageImage, 0)
	err := r.db.SelectContext(ctx, &items, `
		SELECT target, position, src, alt, sort_order FROM blog_page_images WHERE blog_id = $1 ORDER BY sort_order ASC, id ASC
	`, blogID)

	return items, err
}
