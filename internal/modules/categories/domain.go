package categories

import "time"

type Category struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	LogoURL     *string   `db:"logo_url"`
	Description *string   `db:"description"`
	IsActive    bool      `db:"is_active"`
	ParentID    *string   `db:"parent_id"`
	Level       int       `db:"level"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CategoryWithParentName struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	LogoURL     *string   `db:"logo_url"`
	Description *string   `db:"description"`
	IsActive    bool      `db:"is_active"`
	ParentID    *string   `db:"parent_id"`
	ParentName  *string   `db:"parent_name"`
	Level       int       `db:"level"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
