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

func (r *CategoryRepo) SaveCategory(ctx context.Context, request CategoryRequestDTO) error {
	query := "INSERT INTO categories(name, slug, parent_id) VALUES ($1,$2,NULLIF($3, '')::uuid)"
	if _, err := r.db.ExecContext(ctx, query, request.Name, request.Slug, request.ParentId); err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *CategoryRepo) UpdateCategory(ctx context.Context, categoryID string, request CategoryRequestDTO) error {
	query := "UPDATE categories SET name=COALESCE(NULLIF($1, ''), name), slug=COALESCE(NULLIF($2, ''), slug), parent_id=COALESCE($3, parent_id) WHERE id=$4"
	if _, err := r.db.ExecContext(ctx, query, request.Name, request.Slug, request.ParentId, categoryID); err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *CategoryRepo) DeleteCategory(ctx context.Context, categoryID string) error {
	query := "DELETE FROM categories WHERE id=$1"
	if _, err := r.db.ExecContext(ctx, query, categoryID); err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *CategoryRepo) GetCategoryByID(ctx context.Context, categoryID string) (*Category, error) {
	var category Category
	query := "SELECT * FROM categories WHERE id=$1"
	if err := r.db.GetContext(ctx, &category, query, categoryID); err != nil {
		fmt.Println("err", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Category not found")
		}
		return nil, shared.PostgresError(err)
	}
	return &category, nil
}

func (r *CategoryRepo) GetAllCategory(ctx context.Context) ([]Category, error) {
	var categories []Category
	query := "SELECT id, name, slug, parent_id,level, created_at, updated_at FROM categories"
	if err := r.db.SelectContext(ctx, &categories, query); err != nil {
		return nil, shared.PostgresError(err)
	}

	return categories, nil
}
