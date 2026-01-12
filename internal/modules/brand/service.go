package brand

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type BrandService struct {
	logger *zap.Logger
	repo   *BrandRepo
}

func NewBrandService(logger *zap.Logger, repo *BrandRepo) *BrandService {
	return &BrandService{
		logger: logger,
		repo:   repo,
	}
}

func (b *BrandService) CreateBrand(ctx context.Context, request BrandsRequestDTO) (*GenericResponseDTO, error) {
	if err := b.repo.SaveBrand(ctx, request); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Brand Created Successfully",
	}, nil
}

func (b *BrandService) UpdateBrand(ctx context.Context, brandId string, request BrandsRequestDTO) (*GenericResponseDTO, error) {
	if err := b.repo.UpdateBrand(ctx, brandId, request); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Brand Created Successfully",
	}, nil
}

func (b *BrandService) DeleteBrand(ctx context.Context, brandID string) (*GenericResponseDTO, error) {
	if err := b.repo.DeleteBrand(ctx, brandID); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Brand Created Successfully",
	}, nil
}

func (b *BrandService) GetBrandByID(ctx context.Context, brandId string) (*BrandResponse, error) {
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

	return response, nil
}

func (b *BrandService) GetAllBrand(ctx context.Context, page, limit int) (*BrandListResponse, error) {
	offset := (page - 1) * limit
	var brandResponse []BrandResponse

	response, err := b.repo.GetAllBrand(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	for _, data := range response.Brands {
		item := BrandResponse{
			ID:          data.ID,
			Name:        data.Name,
			Slug:        data.Slug,
			LogoUrl:     data.LogoUrl,
			WebsiteUrl:  data.WebsiteUrl,
			Description: data.Description,
			IsActive:    data.IsActive,
			CreatedAt:   data.CreatedAt,
		}

		brandResponse = append(brandResponse, item)
	}

	return &BrandListResponse{
		Brands: brandResponse,
		Total:  response.Total,
	}, nil
}
