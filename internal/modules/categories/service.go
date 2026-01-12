package categories

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type CategoryService struct {
	logger *zap.Logger
	repo   *CategoryRepo
}

func NewCategoryService(logger *zap.Logger, repo *CategoryRepo) *CategoryService {
	return &CategoryService{
		logger: logger,
		repo:   repo,
	}
}

func (b *CategoryService) CreateCategory(ctx context.Context, request CategoryRequestDTO) (*GenericResponseDTO, error) {
	if err := b.repo.SaveCategory(ctx, request); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Category Created Successfully",
	}, nil
}

func (b *CategoryService) UpdateCategory(ctx context.Context, categoryId string, request CategoryRequestDTO) (*GenericResponseDTO, error) {
	if err := b.repo.UpdateCategory(ctx, categoryId, request); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Category Updated Successfully",
	}, nil
}

func (b *CategoryService) DeleteCategory(ctx context.Context, categoryID string) (*GenericResponseDTO, error) {
	if err := b.repo.DeleteCategory(ctx, categoryID); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Category Deleted Successfully",
	}, nil
}

func (b *CategoryService) GetCategoryByID(ctx context.Context, categoryId string) (*CategoryResponse, error) {
	resp, err := b.repo.GetCategoryByID(ctx, categoryId)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	response := &CategoryResponse{
		ID:        resp.ID,
		Name:      resp.Name,
		Slug:      resp.Slug,
		ParentId:  resp.ParentId,
		Level:     resp.Level,
		CreatedAt: resp.CreatedAt,
	}

	return response, nil
}

func (b *CategoryService) GetAllCategory(ctx context.Context) (*GetAllCategory, error) {
	response, err := b.repo.GetAllCategory(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	categories := make([]CategoryResponse, 0, len(response))

	for _, data := range response {
		categories = append(categories, CategoryResponse{
			ID:        data.ID,
			Name:      data.Name,
			Slug:      data.Slug,
			ParentId:  data.ParentId,
			Level:     data.Level,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
		})
	}

	return &GetAllCategory{
		Categories: categories,
	}, nil
}
