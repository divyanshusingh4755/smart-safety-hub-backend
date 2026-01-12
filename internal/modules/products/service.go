package products

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type ProductService struct {
	logger *zap.Logger
	repo   *ProductRepo
}

func NewProductService(logger *zap.Logger, repo *ProductRepo) *ProductService {
	return &ProductService{
		logger: logger,
		repo:   repo,
	}
}

func (b *ProductService) CreateProduct(ctx context.Context, request ProductRequestDTO) (*GenericResponseDTO, error) {
	productId, err := b.repo.SaveProduct(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	return &GenericResponseDTO{
		ID:      productId,
		Status:  "success",
		Message: "Product Created Successfully",
	}, nil
}

func (b *ProductService) UpdateProduct(ctx context.Context, productId string, request ProductRequestDTO) (*GenericResponseDTO, error) {
	if err := b.repo.UpdateProduct(ctx, productId, request); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Product Updated Successfully",
	}, nil
}

func (b *ProductService) DeleteProduct(ctx context.Context, productID string, status string) (*GenericResponseDTO, error) {
	if err := b.repo.DeleteProduct(ctx, productID, status); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Product Archived Successfully",
	}, nil
}

func (b *ProductService) GetProductByID(ctx context.Context, productId string) (*ProductResponseDTO, error) {
	resp, err := b.repo.GetProductByID(ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	response := &ProductResponseDTO{
		ID:          resp.ID,
		Name:        resp.Name,
		Slug:        resp.Slug,
		Description: resp.Description,
		SellerID:    resp.SellerID,
		BrandID:     resp.BrandID,
		CategoryID:  resp.CategoryID,
		Status:      resp.Status,
		CreatedAt:   resp.CreatedAt,
	}

	return response, nil
}

func (b *ProductService) GetProductBySlug(ctx context.Context, slug string) (*ProductResponseDTO, error) {
	resp, err := b.repo.GetProductBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	response := &ProductResponseDTO{
		ID:          resp.ID,
		Name:        resp.Name,
		Slug:        resp.Slug,
		Description: resp.Description,
		SellerID:    resp.SellerID,
		BrandID:     resp.BrandID,
		CategoryID:  resp.CategoryID,
		Status:      resp.Status,
		CreatedAt:   resp.CreatedAt,
	}

	return response, nil
}

func (b *ProductService) GetAllProducts(ctx context.Context, request ProductFilters) (*ProductListResponse, error) {
	resp, err := b.repo.GetAllProducts(ctx, request)

	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	totalCount := 0
	if len(resp) > 0 {
		totalCount = resp[0].TotalCount
	}

	productsDTO := make([]GetProductsDTO, 0, len(resp))
	for _, data := range resp {
		productsDTO = append(productsDTO, GetProductsDTO{
			ID:           data.ID,
			Name:         data.Name,
			Slug:         data.Slug,
			Description:  data.Description,
			Status:       data.Status,
			CategoryName: data.CategoryName,
			BrandName:    data.BrandName,
			ImageURL:     data.ImageURL,
		})
	}
	return &ProductListResponse{
		Products:   productsDTO,
		TotalCount: totalCount,
		Page:       request.Page,
		Limit:      request.Limit,
	}, nil
}

func (b *ProductService) AddProductAttribute(ctx context.Context, request ProductAttributeDTO) (*GenericResponseDTO, error) {
	productAttribute := make([]ProductAttributeArrayDTO, 0, len(request.Attributes))

	for _, data := range request.Attributes {
		productAttribute = append(productAttribute, ProductAttributeArrayDTO{
			ProductID:      request.ProductID,
			AttributeKey:   data.AttributeKey,
			AttributeValue: data.AttributeValue,
		})
	}

	err := b.repo.AddProductAttribute(ctx, request.ProductID, productAttribute)
	if err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	return &GenericResponseDTO{
		ID:      &request.ProductID,
		Status:  "success",
		Message: "Product Attribute Created Successfully",
	}, nil
}

func (b *ProductService) GetProductAttributeByID(ctx context.Context, productId string) (*ProductAttributeDTO, error) {
	response, err := b.repo.GetProductAttributeByID(ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	var productAttribute []ProductAttributeArray

	for _, data := range response {
		productAttribute = append(productAttribute, ProductAttributeArray{
			AttributeKey:   data.AttributeKey,
			AttributeValue: data.AttributeValue,
		})
	}

	return &ProductAttributeDTO{
		ProductID:  productId,
		Attributes: productAttribute,
	}, nil
}

func (b *ProductService) SyncProductVariants(ctx context.Context, productId string, request VariantRequestDTO) (*GenericResponseDTO, error) {
	err := b.repo.SyncProductVariants(ctx, productId, request)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	return &GenericResponseDTO{
		ID:      &productId,
		Status:  "success",
		Message: "Product Variant Saved Successfully",
	}, nil
}

func (b *ProductService) AddProductMedia(ctx context.Context, productId string, request []ProductMediaDTO) (*GenericResponseDTO, error) {
	err := b.repo.AddProductMedia(ctx, productId, request)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	return &GenericResponseDTO{
		ID:      &productId,
		Status:  "success",
		Message: "Product Media Saved Successfully",
	}, nil
}

func (b *ProductService) GetProductMedia(ctx context.Context, productId string) (*[]ProductMediaDTO, error) {
	response, err := b.repo.GetProductMedia(ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	mediaData := make([]ProductMediaDTO, 0, len(response))

	for _, data := range response {
		res := ProductMediaDTO{
			ID:           &data.ID,
			ProductID:    data.ProductID,
			VariantID:    data.VariantID,
			Url:          data.Url,
			MediaType:    data.Type,
			DisplayOrder: data.DisplayOrder,
		}

		mediaData = append(mediaData, res)

	}

	return &mediaData, nil
}

func (b *ProductService) GetProductVariants(ctx context.Context, productId string) (*VariantRequestDTO, error) {
	var response VariantRequestDTO

	result, err := b.repo.GetProductVariants(ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	response.ProductID = result.ProductID
	if err := json.Unmarshal(result.Options, &response.Options); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(result.Variants, &response.Variants); err != nil {
		return nil, err
	}

	return &response, nil
}

func (b *ProductService) SaveProductSEO(ctx context.Context, productId string, request ProductSEODTO) error {
	keywordsJSON, err := json.Marshal(request.Keywords)
	if err != nil {
		return fmt.Errorf("failed to marshal keywords: %v", err)
	}

	seoEntity := ProductSEO{
		ProductID:       productId,
		MetaTitle:       request.MetaTitle,
		MetaDescription: request.MetaDescription,
		OgImageUrl:      request.OgImageUrl,
		Keywords:        json.RawMessage(keywordsJSON),
	}

	return b.repo.SaveProductSEO(ctx, seoEntity)
}

func (b *ProductService) GetProductSEO(ctx context.Context, productId string) (*ProductSEODTO, error) {
	result, err := b.repo.GetProductSEO(ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("failed to get SEO data: %v", err)
	}

	var keywords []string
	if len(result.Keywords) > 0 {
		if err := json.Unmarshal(result.Keywords, &keywords); err != nil {
			return nil, fmt.Errorf("failed to unmarshal keywords: %v", err)
		}
	}

	return &ProductSEODTO{
		ProductID:       result.ProductID,
		MetaTitle:       result.MetaTitle,
		MetaDescription: result.MetaDescription,
		OgImageUrl:      result.OgImageUrl,
		Keywords:        keywords,
	}, nil
}
