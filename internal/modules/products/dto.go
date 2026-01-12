package products

import "time"

type ProductRequestDTO struct {
	Name        string        `json:"name" validate:"required"`
	Slug        string        `json:"slug" validate:"required"`
	Description *string       `json:"description"`
	SellerID    string        `json:"seller_id" validate:"required"`
	BrandID     string        `json:"brand_id" validate:"required"`
	CategoryID  string        `json:"category_id" validate:"required"`
	Status      ProductStatus `json:"status" validate:"required"`
}

type ProductResponseDTO struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Description *string       `json:"description"`
	SellerID    string        `json:"seller_id"`
	BrandID     string        `json:"brand_id"`
	CategoryID  string        `json:"category_id"`
	Status      ProductStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type ProductListResponse struct {
	Products   []GetProductsDTO `json:"products"`
	TotalCount int              `json:"total_count"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
}

type GetProductsDTO struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Slug         string        `json:"slug"`
	Description  *string       `json:"description"`
	BrandName    string        `json:"brand_name"`
	CategoryName string        `json:"category_name"`
	Status       ProductStatus `json:"status"`
	ImageURL     *string       `json:"image_url"`
}

type ProductVariant struct {
	ID           *string  `json:"id"`
	SKU          string   `json:"sku" validate:"required,max=100"`
	Price        float64  `json:"price" validate:"required,gte=0"`
	Weight       float64  `json:"weight" validate:"gte=0"`
	IsActive     bool     `json:"is_active"`
	OptionValues []string `json:"option_values" validate:"required,dive"`
}

type ProductOptionValue struct {
	Name   string   `json:"name" validate:"required"`
	Values []string `json:"values" validate:"required"`
}

type VariantRequestDTO struct {
	ProductID string               `json:"product_id" validate:"required"`
	Options   []ProductOptionValue `json:"options" validate:"required,dive"`
	Variants  []ProductVariant     `json:"variants" validate:"required,dive"`
}

type ProductAttributeDTO struct {
	ProductID  string                  `json:"product_id" validate:"required"`
	Attributes []ProductAttributeArray `json:"attributes" validate:"required"`
}

type ProductAttributeArray struct {
	AttributeKey   string `json:"attribute_key" validate:"required"`
	AttributeValue string `json:"attribute_value" validate:"required"`
}

type ProductAttributeArrayDTO struct {
	ProductID      string `json:"product_id" validate:"required"`
	AttributeKey   string `json:"attribute_key" validate:"required"`
	AttributeValue string `json:"attribute_value" validate:"required"`
}

type ProductMediaDTO struct {
	ID           *string     `json:"id"`
	ProductID    string      `json:"product_id" validate:"required"`
	VariantID    *string     `json:"variant_id"`
	Url          string      `json:"url" validate:"required,url"`
	MediaType    ProductType `json:"type" validate:"required"`
	DisplayOrder int         `json:"display_order" validate:"min=0"`
}

type ProductSEODTO struct {
	ProductID       string   `json:"product_id"`
	MetaTitle       string   `json:"meta_title"`
	MetaDescription string   `json:"meta_description"`
	OgImageUrl      string   `json:"og_image_url"`
	Keywords        []string `json:"keywords"`
}

type ProductFilters struct {
	Category []string `query:"category"`
	Brand    []string `query:"brand"`
	Search   string   `query:"search"`
	Status   string   `query:"status"`
	MinPrice float64  `query:"min_price"`
	MaxPrice float64  `query:"max_price"`
	Page     int      `query:"page"`
	Limit    int      `query:"limit"`
}

type GetProductByID struct {
	ProductID string `json:"id" validate:"required"`
}

type GetProductAttributeByID struct {
	AttributeID string `json:"id" validate:"required"`
}

type GetProductAttributeByProductID struct {
	ProductID string `json:"product_id" validate:"required"`
}

type GenericResponseDTO struct {
	ID      *string `json:"product_id"`
	Status  string  `json:"success"`
	Message string  `json:"message"`
}
