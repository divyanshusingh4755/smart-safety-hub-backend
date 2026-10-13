package blogs

import (
	"encoding/json"
	"time"
)

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
	ID              string                 `json:"id"`
	Title           string                 `json:"title"`
	Slug            string                 `json:"slug"`
	Type            string                 `json:"type"`
	Excerpt         *string                `json:"excerpt"`
	Content         string                 `json:"content"`
	FeaturedImage   *string                `json:"featured_image"`
	MetaTitle       *string                `json:"meta_title"`
	MetaDescription *string                `json:"meta_description"`
	IsPublished     bool                   `json:"is_published"`
	PublishedAt     *time.Time             `json:"published_at"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	PageConfig      *PageConfigResponseDTO `json:"pageConfig",omitempty`
}

type GetAllBlogs struct {
	Blogs []BlogResponse `json:"blogs"`
}

type GenericResponseDTO struct {
	Status  string `json:"success"`
	Message string `json:"message"`
}

type PageConfigRequestDTO struct {
	Layout        string                 `json:"layout" validate:"required,oneof=standard country-landing"`
	Market        json.RawMessage        `json:"market"`
	Hero          json.RawMessage        `json:"hero"`
	Trust         json.RawMessage        `json:"trust"`
	Components    json.RawMessage        `json:"components"`
	Sections      json.RawMessage        `json:"sections"`
	FinalCTA      json.RawMessage        `json:"finalCta"`
	Products      json.RawMessage        `json:"products"`
	Insertions    []BlogPageInsertion    `json:"insertions"`
	InternalLinks []BlogPageInternalLink `json:"internalLinks"`
	Images        []BlogPageImage        `json:"images"`
}

type PageConfigResponseDTO struct {
	BlogID        string                 `json:"blog_id"`
	Layout        string                 `json:"layout"`
	Market        json.RawMessage        `json:"market"`
	Hero          json.RawMessage        `json:"hero"`
	Trust         json.RawMessage        `json:"trust"`
	Components    json.RawMessage        `json:"components"`
	Sections      json.RawMessage        `json:"sections"`
	FinalCTA      json.RawMessage        `json:"finalCta"`
	Products      json.RawMessage        `json:"products"`
	Insertions    []BlogPageInsertion    `json:"insertions"`
	InternalLinks []BlogPageInternalLink `json:"internalLinks"`
	Images        []BlogPageImage        `json:"images"`
}
