package categories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/smart-safety-hub/backend/shared"
)

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{
		db: db,
	}
}

func (r *CategoryRepo) GetLevelByParent(
	ctx context.Context,
	parentID *string,
) (int, error) {

	if parentID == nil {
		return 0, nil
	}

	var level int

	err := r.db.GetContext(
		ctx,
		&level,
		`SELECT level + 1
		 FROM categories
		 WHERE id = $1
		 AND is_active = true`,
		*parentID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("parent category not found")
	}

	return level, err
}

func (r *CategoryRepo) SaveCategory(ctx context.Context, request CategoryRequestDTO) error {
	if request.ParentID != nil && *request.ParentID == "" {
		request.ParentID = nil
	}

	level, err := r.GetLevelByParent(ctx, request.ParentID)
	if err != nil {
		return shared.PostgresError(err)
	}

	query := `
		INSERT INTO categories (
		    name,
		    slug,
		    logo_url,
		    description,
		    is_active,
		    parent_id,
		    level
		)
		VALUES (
		    $1,
		    $2,
		    $3,
		    $4,
		    COALESCE($5, TRUE),
		    NULLIF($6, '')::uuid,
		    $7
		)
		`
	if _, err = r.db.ExecContext(
		ctx,
		query,
		request.Name,
		request.Slug,
		request.LogoURL,
		request.Description,
		request.IsActive,
		request.ParentID,
		level,
	); err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *CategoryRepo) UpdateCategory(ctx context.Context, categoryID string, request UpdateCategoryDTO) error {
	query := `
		UPDATE categories
		SET
		    name = COALESCE($1, name),
		    slug = COALESCE($2, slug),
		    logo_url = COALESCE($3, logo_url),
		    description = COALESCE($4, description),
		    is_active = COALESCE($5, is_active),
		    parent_id = COALESCE(NULLIF($6, '')::uuid, parent_id),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		`
	if _, err := r.db.ExecContext(
		ctx,
		query,
		request.Name,
		request.Slug,
		request.LogoURL,
		request.Description,
		request.IsActive,
		request.ParentID,
		categoryID,
	); err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *CategoryRepo) DeleteCategory(ctx context.Context, categoryID string) error {
	query := `
	UPDATE categories
	SET
		is_active = false,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, categoryID)
	if err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *CategoryRepo) GetCategoryByID(ctx context.Context, categoryID string) (*Category, error) {
	var category Category
	query := `SELECT *
		FROM categories
		WHERE id=$1`
	if err := r.db.GetContext(ctx, &category, query, categoryID); err != nil {
		fmt.Printf("SQL ERROR: %+v\n", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Category not found")
		}
		return nil, shared.PostgresError(err)
	}
	return &category, nil
}

func (r *CategoryRepo) GetCategoryBySlug(ctx context.Context, categorySlug string) (*Category, error) {
	var category Category
	query := `SELECT *
		FROM categories
		WHERE slug=$1
		AND is_active = true`
	if err := r.db.GetContext(ctx, &category, query, categorySlug); err != nil {
		fmt.Println("err", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Category not found")
		}
		return nil, shared.PostgresError(err)
	}
	return &category, nil
}

func (r *CategoryRepo) GetAllCategory(ctx context.Context) ([]CategoryWithParentName, error) {
	var categories []CategoryWithParentName
	query := `
	        SELECT c.*, p.name AS parent_name
	        FROM categories c
	        LEFT JOIN categories p ON c.parent_id = p.id`
	if err := r.db.SelectContext(ctx, &categories, query); err != nil {
		return nil, shared.PostgresError(err)
	}

	return categories, nil
}
