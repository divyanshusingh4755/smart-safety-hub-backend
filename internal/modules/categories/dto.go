package categories

import "time"

type CategoryRequestDTO struct {
	Name     string  `json:"name" validate:"required"`
	Slug     string  `json:"slug" validate:"required"`
	ParentId *string `json:"parent_id"`
}

type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ParentId  *string   `json:"parent_id"`
	Level     *int      `json:"level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetAllCategory struct {
	Categories []CategoryResponse `json:"categories"`
}

type GetCategoryByID struct {
	CategoryID string `json:"id" validate:"required"`
}

type GenericResponseDTO struct {
	Status  string `json:"success"`
	Message string `json:"message"`
}
