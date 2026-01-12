package products

import (
	"encoding/json"
	"time"
)

type ProductStatus string

type ProductType string

const (
	DRAFT    ProductStatus = "DRAFT"
	ACTIVE   ProductStatus = "ACTIVE"
	ARCHIVED ProductStatus = "ARCHIVED"
)

const (
	IMAGE ProductType = "image"
	VIDEO ProductType = "video"
	PDF   ProductType = "pdf"
)

type Product struct {
	ID          string        `db:"id"`
	Name        string        `db:"name"`
	Slug        string        `db:"slug"`
	Description *string       `db:"description"`
	SellerID    string        `db:"seller_id"`
	BrandID     string        `db:"brand_id"`
	CategoryID  string        `db:"category_id"`
	Status      ProductStatus `db:"status"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
}

type GetProducts struct {
	ID           string        `db:"id"`
	Name         string        `db:"name"`
	Slug         string        `db:"slug"`
	Description  *string       `db:"description"`
	BrandName    string        `db:"brand_name"`
	CategoryName string        `db:"category_name"`
	Status       ProductStatus `db:"status"`
	ImageURL     *string       `db:"image_url"`
	TotalCount   int           `db:"total_count"`
}

type ProductAttribute struct {
	ID             string    `db:"id"`
	ProductID      string    `db:"product_id"`
	AttributeKey   string    `db:"attribute_key"`
	AttributeValue string    `db:"attribute_value"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type ProductOptions struct {
	ID        string    `db:"id"`
	ProductID string    `db:"product_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ProductOptionValues struct {
	ID        string    `db:"id"`
	OptionID  string    `db:"option_id"`
	Value     string    `db:"value"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ProductVariants struct {
	ProductID string          `db:"product_id"`
	Options   json.RawMessage `db:"options"`
	Variants  json.RawMessage `db:"variants"`
}

type ProductMedia struct {
	ID           string      `db:"id"`
	ProductID    string      `db:"product_id"`
	VariantID    *string     `db:"variant_id"`
	Url          string      `db:"url"`
	Type         ProductType `db:"type"`
	DisplayOrder int         `db:"display_order"`
}

type ProductSEO struct {
	ProductID       string          `db:"product_id"`
	MetaTitle       string          `db:"meta_title"`
	MetaDescription string          `db:"meta_description"`
	OgImageUrl      string          `db:"og_image_url"`
	Keywords        json.RawMessage `db:"keywords"`
}
