package brand

import "time"

type Brand struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	LogoUrl     *string   `db:"logo_url"`
	WebsiteUrl  *string   `db:"website_url"`
	Description *string   `db:"description"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	TotalCount  int       `db:"total_count"`
}

type BrandList struct {
	Brands []Brand `db:"brands"`
	Total  int     `db:"total"`
}
