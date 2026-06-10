package brand

import (
	"context"
	"fmt"

	"github.com/smart-safety-hub/backend/shared/cache"
	"go.uber.org/zap"
)

type BrandService struct {
	logger *zap.Logger
	repo   *BrandRepo
	cache  *cache.RedisCache
}

func NewBrandService(logger *zap.Logger, repo *BrandRepo, cache *cache.RedisCache) *BrandService {
	return &BrandService{
		logger: logger,
		repo:   repo,
		cache:  cache,
	}
}

func (b *BrandService) CreateBrand(ctx context.Context, request BrandsRequestDTO) (*GenericResponseDTO, error) {
	if err := b.repo.SaveBrand(ctx, request); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	b.invalidateBrandCache(ctx)

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Brand Created Successfully",
	}, nil
}

func (b *BrandService) UpdateBrand(ctx context.Context, brandId string, request BrandsRequestDTO) (*GenericResponseDTO, error) {
	if err := b.repo.UpdateBrand(ctx, brandId, request); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	b.invalidateBrandCache(ctx)

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Brand Created Successfully",
	}, nil
}

func (b *BrandService) DeleteBrand(ctx context.Context, brandID string) (*GenericResponseDTO, error) {
	if err := b.repo.DeleteBrand(ctx, brandID); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	b.invalidateBrandCache(ctx)

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Brand Created Successfully",
	}, nil
}

func (b *BrandService) GetBrandByID(ctx context.Context, brandId string) (*BrandResponse, error) {
	cacheKey := cache.BrandByID(brandId)

	var cached BrandResponse
	if err := b.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}
	resp, err := b.repo.GetBrandByID(ctx, brandId)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	response := &BrandResponse{
		ID:          resp.ID,
		Name:        resp.Name,
		Slug:        resp.Slug,
		LogoUrl:     resp.LogoUrl,
		WebsiteUrl:  resp.WebsiteUrl,
		Description: resp.Description,
		IsActive:    resp.IsActive,
		CreatedAt:   resp.CreatedAt,
	}

	_ = b.cache.Set(ctx, cacheKey, response, cache.BrandTTL)

	return response, nil
}

func (b *BrandService) GetBrandBySlug(ctx context.Context, brandSlug string) (*BrandResponse, error) {
	cacheKey := cache.BrandByID(brandSlug)

	var cached BrandResponse
	if err := b.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}
	resp, err := b.repo.GetBrandBySlug(ctx, brandSlug)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	response := &BrandResponse{
		ID:          resp.ID,
		Name:        resp.Name,
		Slug:        resp.Slug,
		LogoUrl:     resp.LogoUrl,
		WebsiteUrl:  resp.WebsiteUrl,
		Description: resp.Description,
		IsActive:    resp.IsActive,
		CreatedAt:   resp.CreatedAt,
	}

	_ = b.cache.Set(ctx, cacheKey, response, cache.BrandTTL)

	return response, nil
}

func (b *BrandService) GetAllBrand(ctx context.Context, page, limit int) (*BrandListResponse, error) {

	cacheKey := cache.BrandList(page, limit)
	var cached BrandListResponse
	if err := b.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	offset := (page - 1) * limit

	response, err := b.repo.GetAllBrand(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	brands := make([]BrandResponse, 0, len(response.Brands))

	for _, data := range response.Brands {
		brands = append(brands, BrandResponse{
			ID:          data.ID,
			Name:        data.Name,
			Slug:        data.Slug,
			LogoUrl:     data.LogoUrl,
			WebsiteUrl:  data.WebsiteUrl,
			Description: data.Description,
			IsActive:    data.IsActive,
			CreatedAt:   data.CreatedAt,
		})
	}

	result := &BrandListResponse{
		Brands: brands,
		Total:  response.Total,
	}

	_ = b.cache.Set(ctx, cacheKey, result, cache.BrandTTL)

	return result, nil
}

func (b *BrandService) invalidateBrandCache(ctx context.Context) {
	keys, err := b.cache.Scan(ctx, "brand:*")
	if err != nil {
		b.logger.Warn("failed to scan cache keys", zap.Error(err))
		return
	}

	if len(keys) > 0 {
		_ = b.cache.Delete(ctx, keys...)
	}
}
