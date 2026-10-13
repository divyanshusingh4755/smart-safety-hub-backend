package blogs

import (
	"encoding/json"
	"time"
)

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

type BlogPageConfig struct {
	ID         string          `db:"id"`
	BlogID     string          `db:"blog_id"`
	Layout     string          `db:"layout"`
	Market     json.RawMessage `db:"market"`
	Hero       json.RawMessage `db:"hero"`
	Trust      json.RawMessage `db:"trust"`
	Components json.RawMessage `db:"components"`
	Sections   json.RawMessage `db:"sections"`
	FinalCTA   json.RawMessage `db:"final_cta"`
	Products   json.RawMessage `db:"products"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
}

type BlogPageInsertion struct {
	Component string `db:"component" json:"component"`
	Target    string `db:"target" json:"target"`
	Position  string `db:"position" json:"position"`
	SortOrder int    `db:"sort_order" json:"-"`
}

type BlogPageInternalLink struct {
	Text           string `db:"text" json:"text"`
	Href           string `db:"href" json:"href"`
	MaxOccurrences int    `db:"max_occurrences" json:"maxOccurrences"`
	SortOrder      int    `db:"sort_order" json:"-"`
}

type BlogPageImage struct {
	Target    string `db:"target" json:"target"`
	Position  string `db:"position" json:"position"`
	Src       string `db:"src" json:"src"`
	Alt       string `db:"alt" json:"alt"`
	SortOrder int    `db:"sort_order" json:"-"`
}
