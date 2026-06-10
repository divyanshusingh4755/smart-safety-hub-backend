package categories

import (
	"context"
	"errors"
	"fmt"

	"github.com/smart-safety-hub/backend/shared/cache"
	"go.uber.org/zap"
)

type CategoryService struct {
	logger *zap.Logger
	repo   *CategoryRepo
	cache  *cache.RedisCache
}

func NewCategoryService(logger *zap.Logger, repo *CategoryRepo, cache *cache.RedisCache) *CategoryService {
	return &CategoryService{
		logger: logger,
		repo:   repo,
		cache:  cache,
	}
}

func (b *CategoryService) CreateCategory(ctx context.Context, request CategoryRequestDTO) (*GenericResponseDTO, error) {
	if err := b.repo.SaveCategory(ctx, request); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	b.invalidateCategoryCache(ctx)

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Category Created Successfully",
	}, nil
}

func (b *CategoryService) UpdateCategory(ctx context.Context, categoryId string, request UpdateCategoryDTO) (*GenericResponseDTO, error) {

	if request.ParentID != nil && *request.ParentID == categoryId {
		return nil, errors.New("category cannot be its own parent")
	}

	if err := b.repo.UpdateCategory(ctx, categoryId, request); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	b.invalidateCategoryCache(ctx)

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Category Updated Successfully",
	}, nil
}

func (b *CategoryService) DeleteCategory(ctx context.Context, categoryID string) (*GenericResponseDTO, error) {
	if err := b.repo.DeleteCategory(ctx, categoryID); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	b.invalidateCategoryCache(ctx)

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Category Deleted Successfully",
	}, nil
}

func (b *CategoryService) GetCategoryByID(ctx context.Context, categoryId string) (*CategoryResponse, error) {

	cacheKey := cache.CategoryByID(categoryId)

	var cached CategoryResponse
	if err := b.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	resp, err := b.repo.GetCategoryByID(ctx, categoryId)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	response := &CategoryResponse{
		ID:          resp.ID,
		Name:        resp.Name,
		Slug:        resp.Slug,
		LogoURL:     resp.LogoURL,
		Description: resp.Description,
		IsActive:    resp.IsActive,
		ParentID:    resp.ParentID,
		Level:       resp.Level,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}

	_ = b.cache.Set(ctx, cacheKey, response, cache.CategoryTTL)

	return response, nil
}

func (b *CategoryService) GetCategoryBySlug(ctx context.Context, categorySlug string) (*CategoryResponse, error) {

	cacheKey := cache.CategoryBySlug(categorySlug)

	var cached CategoryResponse
	if err := b.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	category, err := b.repo.GetCategoryBySlug(ctx, categorySlug)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	response := &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		LogoURL:     category.LogoURL,
		Description: category.Description,
		IsActive:    category.IsActive,
		ParentID:    category.ParentID,
		Level:       category.Level,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	_ = b.cache.Set(ctx, cacheKey, response, cache.CategoryTTL)

	return response, nil
}

func (s *CategoryService) GetAllCategory(ctx context.Context) (*GetAllCategory, error) {
	cacheKey := cache.CategoryList()

	var cached GetAllCategory
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	categories, err := s.repo.GetAllCategory(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching categories from DB: %v", err)
	}

	respList := make([]CategoryWithParentNameResponse, len(categories))
	for i, cat := range categories {
		var logoURL, description string

		if cat.LogoURL != nil {
			logoURL = *cat.LogoURL
		}

		if cat.Description != nil {
			description = *cat.Description
		}

		respList[i] = CategoryWithParentNameResponse{
			ID:          cat.ID,
			Name:        cat.Name,
			Slug:        cat.Slug,
			LogoURL:     logoURL,
			Description: description,
			IsActive:    cat.IsActive,
			ParentID:    cat.ParentID,
			ParentName:  cat.ParentName,
			Level:       cat.Level,
			CreatedAt:   cat.CreatedAt,
			UpdatedAt:   cat.UpdatedAt,
		}
	}

	result := &GetAllCategory{Categories: respList}
	_ = s.cache.Set(ctx, cacheKey, result, cache.CategoryTTL)

	return result, nil
}

func (b *CategoryService) invalidateCategoryCache(ctx context.Context) {
	keys, err := b.cache.Scan(ctx, "category:*")
	if err != nil {
		b.logger.Warn("failed to scan cache keys", zap.Error(err))
		return
	}

	if len(keys) > 0 {
		_ = b.cache.Delete(ctx, keys...)
	}
}
