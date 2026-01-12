package brand

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/smart-safety-hub/backend/shared"
)

type BrandRepo struct {
	db *sqlx.DB
}

func NewBrandRepo(db *sqlx.DB) *BrandRepo {
	return &BrandRepo{
		db: db,
	}
}

func (r *BrandRepo) SaveBrand(ctx context.Context, request BrandsRequestDTO) error {
	query := "INSERT INTO brands(name, slug, logo_url, website_url, description) VALUES ($1,$2,$3,$4,$5)"
	if _, err := r.db.ExecContext(ctx, query, request.Name, request.Slug, request.LogoUrl, request.WebsiteUrl, request.Description); err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *BrandRepo) UpdateBrand(ctx context.Context, brandID string, request BrandsRequestDTO) error {
	query := "UPDATE brands SET name=COALESCE(NULLIF($1, ''), name), slug=COALESCE(NULLIF($2, ''), slug), logo_url=COALESCE($3, logo_url), website_url=COALESCE($4, website_url), description=COALESCE($5, description) WHERE id=$6"
	if _, err := r.db.ExecContext(ctx, query, request.Name, request.Slug, request.LogoUrl, request.WebsiteUrl, request.Description, brandID); err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *BrandRepo) DeleteBrand(ctx context.Context, brandID string) error {
	query := "DELETE FROM brands WHERE id=$1"
	if _, err := r.db.ExecContext(ctx, query, brandID); err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *BrandRepo) GetBrandByID(ctx context.Context, brandId string) (*Brand, error) {
	var brand Brand
	query := "SELECT * FROM brands WHERE id=$1"
	if err := r.db.GetContext(ctx, &brand, query, brandId); err != nil {
		fmt.Println("err", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Brand not found")
		}
		return nil, shared.PostgresError(err)
	}
	return &brand, nil
}

func (r *BrandRepo) GetAllBrand(limit, offset int) (*BrandList, error) {
	var brands []Brand
	query := "SELECT id, name, slug, logo_url, is_active, website_url, created_at, COUNT(*) OVER() as total_count FROM brands ORDER BY id DESC LIMIT $1 OFFSET $2"
	if err := r.db.Select(&brands, query, limit, offset); err != nil {
		return nil, shared.PostgresError(err)
	}

	total := 0
	if len(brands) > 0 {
		total = brands[0].TotalCount
	}

	return &BrandList{
		Brands: brands,
		Total:  total,
	}, nil
}
