package blogs

import "time"

type BlogRequestDTO struct {
	Title           string  `json:"title" validate:"required"`
	Slug            string  `json:"slug" validate:"required"`
	Type            string  `json:"type" validate:"required,oneof=post page"`
	Excerpt         *string `json:"excerpt"`
	Content         string  `json:"content" validate:"required"`
	FeaturedImage   *string `json:"featured_image"`
	MetaTitle       *string `json:"meta_title"`
	MetaDescription *string `json:"meta_description"`
	IsPublished     *bool   `json:"is_published"`
}

type UpdateBlogDTO struct {
	Title           *string `json:"title"`
	Slug            *string `json:"slug"`
	Type            *string `json:"type" validate:"omitempty,oneof=post page"`
	Excerpt         *string `json:"excerpt"`
	Content         *string `json:"content"`
	FeaturedImage   *string `json:"featured_image"`
	MetaTitle       *string `json:"meta_title"`
	MetaDescription *string `json:"meta_description"`
	IsPublished     *bool   `json:"is_published"`
}

type BlogResponse struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	Type            string     `json:"type"`
	Excerpt         *string    `json:"excerpt"`
	Content         string     `json:"content"`
	FeaturedImage   *string    `json:"featured_image"`
	MetaTitle       *string    `json:"meta_title"`
	MetaDescription *string    `json:"meta_description"`
	IsPublished     bool       `json:"is_published"`
	PublishedAt     *time.Time `json:"published_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type GetAllBlogs struct {
	Blogs []BlogResponse `json:"blogs"`
}

type GenericResponseDTO struct {
	Status  string `json:"success"`
	Message string `json:"message"`
}
