package products

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/smart-safety-hub/backend/shared/cache"
	"go.uber.org/zap"
)

type ProductService struct {
	logger *zap.Logger
	repo   *ProductRepo
	cache  *cache.RedisCache
}

func NewProductService(logger *zap.Logger, repo *ProductRepo, cache *cache.RedisCache) *ProductService {
	return &ProductService{
		logger: logger,
		repo:   repo,
		cache:  cache,
	}
}

func (b *ProductService) CreateProduct(ctx context.Context, request ProductRequestDTO) (*GenericResponseDTO, error) {
	productId, err := b.repo.SaveProduct(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	b.invalidateListCache(ctx)

	return &GenericResponseDTO{
		ID:      productId,
		Status:  "success",
		Message: "Product Created Successfully",
	}, nil
}

func (b *ProductService) UpdateProduct(ctx context.Context, productId string, request ProductRequestDTO) (*GenericResponseDTO, error) {
	old, err := b.repo.GetProductByID(ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch old product: %w", err)
	}

	if err := b.repo.UpdateProduct(ctx, productId, request); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	if old != nil {
		b.inValidateProductCache(ctx, productId, old.Slug)
	}

	b.invalidateListCache(ctx)

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Product Updated Successfully",
	}, nil
}

func (b *ProductService) DeleteProduct(ctx context.Context, productID string, status string) (*GenericResponseDTO, error) {
	if err := b.repo.DeleteProduct(ctx, productID, status); err != nil {
		return nil, fmt.Errorf("Error came while saving it to DB: %v", err)
	}

	product, _ := b.repo.GetProductByID(ctx, productID)

	if product != nil {
		b.inValidateProductCache(ctx, productID, product.Slug)
	}
	b.invalidateListCache(ctx)

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Product Archived Successfully",
	}, nil
}

func (b *ProductService) GetProductByID(ctx context.Context, productId string) (*ProductResponseDTO, error) {

	cacheKey := cache.ProductById(productId)
	var cached ProductResponseDTO
	if err := b.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	resp, err := b.repo.GetProductByID(ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %w", err)
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

	_ = b.cache.Set(ctx, cacheKey, response, cache.ProductTTL)
	return response, nil
}

func (b *ProductService) GetProductBySlug(ctx context.Context, slug string) (*ProductResponseDTO, error) {

	cacheKey := cache.ProductBySlug(slug)

	var cached ProductResponseDTO

	if err := b.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	resp, err := b.repo.GetProductBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	response := &ProductResponseDTO{
		ID:           resp.ID,
		Name:         resp.Name,
		Slug:         resp.Slug,
		Description:  resp.Description,
		SellerID:     resp.SellerID,
		BrandID:      resp.BrandID,
		CategoryID:   resp.CategoryID,
		Status:       resp.Status,
		CreatedAt:    resp.CreatedAt,
		UpdatedAt:    resp.UpdatedAt,
		BrandName:    resp.BrandName,
		BrandSlug:    resp.BrandSlug,
		CategoryName: resp.CategoryName,
		CategorySlug: resp.CategorySlug,
	}

	_ = b.cache.Set(ctx, cacheKey, response, cache.ProductTTL)
	return response, nil
}

func (b *ProductService) GetAllProducts(ctx context.Context, request ProductFilters) (*ProductListResponse, error) {
	version := b.cache.GetVersion(ctx)
	cacheKey := buildProductListCacheKey(request, version)

	var cached ProductListResponse
	if err := b.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	resp, err := b.repo.GetAllProducts(ctx, request)
	fmt.Print("res", resp)
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

	response := &ProductListResponse{
		Products:   productsDTO,
		TotalCount: totalCount,
		Page:       request.Page,
		Limit:      request.Limit,
	}

	_ = b.cache.Set(ctx, cacheKey, response, cache.ProductTTL)
	return response, nil
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

	b.inValidateProductCache(ctx, request.ProductID, "")

	return &GenericResponseDTO{
		ID:      &request.ProductID,
		Status:  "success",
		Message: "Product Attribute Created Successfully",
	}, nil
}

func (b *ProductService) GetProductAttributeByID(ctx context.Context, productId string) (*ProductAttributeDTO, error) {

	cacheKey := cache.ProductAttributes(productId)
	var cached ProductAttributeDTO

	if err := b.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	response, err := b.repo.GetProductAttributeByID(ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	attrs := make([]ProductAttributeArray, 0, len(response))
	for _, data := range response {
		attrs = append(attrs, ProductAttributeArray{
			AttributeKey:   data.AttributeKey,
			AttributeValue: data.AttributeValue,
		})
	}

	dto := &ProductAttributeDTO{
		ProductID:  productId,
		Attributes: attrs,
	}

	_ = b.cache.Set(ctx, cacheKey, dto, cache.ProductTTL)

	return dto, nil
}

func (b *ProductService) SyncProductVariants(ctx context.Context, productId string, request VariantRequestDTO) (*GenericResponseDTO, error) {
	err := b.repo.SyncProductVariants(ctx, productId, request)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting data from DB: %v", err)
	}

	b.inValidateProductCache(ctx, productId, "")

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

	b.inValidateProductCache(ctx, productId, "")

	return &GenericResponseDTO{
		ID:      &productId,
		Status:  "success",
		Message: "Product Media Saved Successfully",
	}, nil
}

func (b *ProductService) GetProductMedia(ctx context.Context, productId string) ([]ProductMediaDTO, error) {

	cacheKey := cache.ProductMedia(productId)
	var cached []ProductMediaDTO

	if err := b.cache.Get(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

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

	_ = b.cache.Set(ctx, cacheKey, mediaData, cache.ProductTTL)

	return mediaData, nil
}

func (b *ProductService) GetProductVariants(ctx context.Context, productId string) (*VariantRequestDTO, error) {
	cacheKey := cache.ProductVariants(productId)
	var response VariantRequestDTO

	if err := b.cache.Get(ctx, cacheKey, &response); err == nil {
		return &response, nil
	}

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

	_ = b.cache.Set(ctx, cacheKey, response, cache.ProductTTL)

	return &response, nil
}

func (b *ProductService) SaveProductSEO(ctx context.Context, productId string, request ProductSEODTO, isPublish ...bool) error {
	shouldPublish := true
	if len(isPublish) > 0 {
		shouldPublish = isPublish[0]
	}

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

	err = b.repo.SaveProductSEO(ctx, seoEntity, shouldPublish)
	if err != nil {
		return err
	}

	b.inValidateProductCache(ctx, productId, "")
	return nil
}

func (b *ProductService) GetProductSEO(ctx context.Context, productId string) (*ProductSEODTO, error) {
	cacheKey := cache.ProductSEO(productId)
	var cached ProductSEODTO

	if err := b.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

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

	dto := &ProductSEODTO{
		ProductID:       result.ProductID,
		MetaTitle:       result.MetaTitle,
		MetaDescription: result.MetaDescription,
		OgImageUrl:      result.OgImageUrl,
		Keywords:        keywords,
	}

	_ = b.cache.Set(ctx, cacheKey, dto, cache.ProductTTL)

	return dto, nil
}

func (b *ProductService) ImportProduct(ctx context.Context, req ImportProductDTO) (string, error) {

	request := ProductRequestDTO{
		Name:        req.ProductName,
		Slug:        req.Slug,
		Description: &req.ProductDescription,
		SellerID:    &req.SellerID,
		BrandID:     nil,
		CategoryID:  nil,
		Status:      "DRAFT",
	}

	productId, err := b.repo.SaveProduct(ctx, request)

	if err != nil {
		return "", fmt.Errorf("Error while saving product to DB: %w", err)
	}

	attribute := ProductAttributeDTO{
		ProductID:  *productId,
		Attributes: req.Attribute,
	}

	productSEO := ProductSEODTO{
		ProductID:       *productId,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		OgImageUrl:      "",
		Keywords:        req.Keywords,
	}

	if _, err = b.AddProductAttribute(ctx, attribute); err != nil {
		return "", fmt.Errorf("Error came while saving product Attribute: %w", err)
	}

	if err = b.SaveProductSEO(ctx, *productId, productSEO, false); err != nil {
		return "", fmt.Errorf("Error came while saving product SEO: %w", err)
	}

	b.invalidateListCache(ctx)
	return "Product Saved Successfully", nil
}

func (b *ProductService) inValidateProductCache(ctx context.Context, ProductById string, slug string) {
	keys := []string{
		cache.ProductById(ProductById),
		cache.ProductBySlug(slug),
		cache.ProductAttributes(ProductById),
		cache.ProductMedia(ProductById),
		cache.ProductVariants(ProductById),
		cache.ProductSEO(ProductById),
	}

	_ = b.cache.Delete(ctx, keys...)
}

func (b *ProductService) invalidateListCache(ctx context.Context) {
	if _, err := b.cache.Incr(ctx, cache.ProductListVersionKey); err != nil {
		b.logger.Warn("failed to increment product list version", zap.Error(err))
	}
}

func buildProductListCacheKey(r ProductFilters, version int64) string {
	brand := ""
	if len(r.Brand) > 0 {
		brand = r.Brand[0]
	}

	category := ""
	if len(r.Category) > 0 {
		category = r.Category[0]
	}

	search := url.QueryEscape(r.Search)

	return fmt.Sprintf(
		"products:list:v%d:p:%d:l:%d:s:%s:b:%s:c:%s:q:%s",
		version,
		r.Page,
		r.Limit,
		r.Status,
		brand,
		category,
		search,
	)
}
