package categories

import "time"

type CategoryRequestDTO struct {
	Name        string  `json:"name" validate:"required"`
	Slug        string  `json:"slug" validate:"required"`
	LogoURL     *string `json:"logo_url"`
	Description *string `json:"description"`
	ParentID    *string `json:"parent_id"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateCategoryDTO struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	LogoURL     *string `json:"logo_url"`
	Description *string `json:"description"`
	ParentID    *string `json:"parent_id"`
	IsActive    *bool   `json:"is_active"`
}

type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	LogoURL     *string   `json:"logo_url"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	ParentID    *string   `json:"parent_id"`
	Level       int       `json:"level"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetAllCategory struct {
	Categories []CategoryWithParentNameResponse `json:"categories"`
}

type GetCategoryByID struct {
	CategoryID string `json:"id" validate:"required"`
}

type GenericResponseDTO struct {
	Status  string `json:"success"`
	Message string `json:"message"`
}

type CategoryWithParentNameResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	LogoURL     string    `json:"logo_url"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	ParentID    *string   `json:"parent_id"`
	ParentName  *string   `json:"parent_name"`
	Level       int       `json:"level"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
