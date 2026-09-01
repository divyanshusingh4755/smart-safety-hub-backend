package quotes

import "time"

const (
	QuoteStatusNew       = "NEW"
	QuoteStatusReviewing = "REVIEWING"
	QuoteStatusQuoted    = "QUOTED"
	QuoteStatusAccepted  = "ACCEPTED"
	QuoteStatusRejected  = "REJECTED"
	QuoteStatusClosed    = "CLOSED"
)

type CreateQuoteItemDTO struct {
	ProductRef string  `json:"product_ref" validate:"required,max=255"`
	Name       string  `json:"name" validate:"omitempty,min=1,max=500"`
	SKU        *string `json:"sku" validate:"omitempty,max=255"`
	Image      *string `json:"image"`
	Quantity   int     `json:"quantity" validate:"required,min=1,max=1000000"`
}

type CreateQuoteRequestDTO struct {
	FullName               string               `json:"full_name" validate:"required,min=2,max=100"`
	CompanyName            string               `json:"company_name" validate:"required,min=2,max=255"`
	Email                  string               `json:"email" validate:"required,max=255"`
	Phone                  string               `json:"phone" validate:"required,max=50"`
	DeliveryLocation       string               `json:"delivery_location" validate:"required,max=255"`
	RequiredBy             *string              `json:"required_by" validate:"omitempty,datetime=2006-01-02"`
	AdditionalRequirements *string              `json:"additional_requirements"`
	Items                  []CreateQuoteItemDTO `json:"items" validate:"required,min=1,dive"`
}

type UpdateQuoteStatusDTO struct {
	Status string `json:"status" validate:"required,oneof=NEW REVIEWING QUOTED ACCEPTED REJECTED CLOSED"`
}

type QuoteItemResponseDTO struct {
	ID         string  `json:"id"`
	ProductRef string  `json:"product_ref"`
	Name       string  `json:"name"`
	SKU        *string `json:"sku"`
	Image      *string `json:"image"`
	Quantity   int     `json:"quantity"`
}

type QuoteResponseDTO struct {
	ID                     string                 `json:"id"`
	QuoteNumber            string                 `json:"quote_number"`
	FullName               string                 `json:"full_name"`
	CompanyName            string                 `json:"company_name"`
	Email                  string                 `json:"email"`
	Phone                  string                 `json:"phone"`
	DeliveryLocation       string                 `json:"delivery_location"`
	RequiredBy             *time.Time             `json:"required_by"`
	AdditionalRequirements *string                `json:"additional_requirements"`
	Status                 string                 `json:"status"`
	Items                  []QuoteItemResponseDTO `json:"items"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
}

type CreateQuoteResponseDTO struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	ID          string `json:"id"`
	QuoteNumber string `json:"quote_number"`
}

type GenericResponseDTO struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GetQuotesQueryDTO struct {
	Page   int    `validate:"min=1"`
	Limit  int    `validate:"min=1,max=100"`
	Search string `validate:"max=255"`
	Status string `validate:"omitempty,oneof=NEW REVIEWING QUOTED ACCEPTED REJECTED CLOSED"`
}

type QuotePaginationDTO struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type GetAllQuotesResponseDTO struct {
	Quotes     []QuoteResponseDTO `json:"quotes"`
	Pagination QuotePaginationDTO `json:"pagination"`
}
