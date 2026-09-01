package quotes

import "time"

type Quote struct {
	ID                     string     `db:"id"`
	QuoteNumber            string     `db:"quote_number"`
	FullName               string     `db:"full_name"`
	CompanyName            string     `db:"company_name"`
	Email                  string     `db:"email"`
	Phone                  string     `db:"phone"`
	DeliveryLocation       string     `db:"delivery_location"`
	RequiredBy             *time.Time `db:"required_by"`
	AdditionalRequirements *string    `db:"additional_requirements"`
	Status                 string     `db:"status"`
	CreatedAt              time.Time  `db:"created_at"`
	UpdatedAt              time.Time  `db:"updated_at"`
}

type QuoteItem struct {
	ID             string    `db:"id"`
	QuoteRequestID string    `db:"quote_request_id"`
	ProductRef     string    `db:"product_ref"`
	Name           string    `db:"product_name"`
	SKU            *string   `db:"sku"`
	Image          *string   `db:"image_url"`
	Quantity       int       `db:"quantity"`
	CreatedAt      time.Time `db:"created_at"`
}
