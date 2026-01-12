package categories

import "time"

type Category struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Slug      string    `db:"slug"`
	ParentId  *string   `db:"parent_id"`
	Level     *int      `db:"level"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
