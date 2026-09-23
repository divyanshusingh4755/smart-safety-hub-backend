package blogs

import "time"

type Blog struct {
	ID              string     `db:"id"`
	Title           string     `db:"title"`
	Slug            string     `db:"slug"`
	Type            string     `db:"type"`
	Excerpt         *string    `db:"excerpt"`
	Content         string     `db:"content"`
	FeaturedImage   *string    `db:"featured_image"`
	MetaTitle       *string    `db:"meta_title"`
	MetaDescription *string    `db:"meta_description"`
	IsPublished     bool       `db:"is_published"`
	PublishedAt     *time.Time `db:"published_at"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}
