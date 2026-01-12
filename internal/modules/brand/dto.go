package brand

import "time"

type BrandsRequestDTO struct {
	Name        string  `json:"name" validate:"required"`
	Slug        string  `json:"slug" validate:"required"`
	LogoUrl     *string `json:"logo_url"`
	WebsiteUrl  *string `json:"website_url"`
	Description *string `json:"description"`
}

type BrandResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	LogoUrl     *string   `json:"logo_url"`
	WebsiteUrl  *string   `json:"website_url"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BrandListResponse struct {
	Brands []BrandResponse `json:"brands"`
	Total  int             `json:"total"`
}

type GetBrandByID struct {
	BrandID string `json:"id" validate:"required"`
}

type GenericResponseDTO struct {
	Status  string `json:"success"`
	Message string `json:"message"`
}
