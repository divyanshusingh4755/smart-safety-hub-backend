package contacts

import "time"

const (
	ContactSourceContactPage = "CONTACT_PAGE"
	ContactSourceFooter      = "FOOTER"
	ContactSourcePopup       = "POPUP"
)

const (
	ContactStatusNew       = "NEW"
	ContactStatusContacted = "CONTACTED"
	ContactStatusQualified = "QUALIFIED"
	ContactStatusClosed    = "CLOSED"
	ContactStatusSpam      = "SPAM"
)

type CreateContactRequestDTO struct {
	Name        string  `json:"name" validate:"required,min=2,max=150"`
	CompanyName *string `json:"company_name"`
	Email       *string `json:"email" validate:"omitempty,email"`
	Phone       string  `json:"phone" validate:"required,max=50"`
	Country     *string `json:"country"`
	InquiryType *string `json:"inquiry_type"`
	ProductName *string `json:"product_name"`
	Quantity    *string `json:"quantity"`
	Message     *string `json:"message"`
	Source      string  `json:"source" validate:"required,oneof=CONTACT_PAGE FOOTER POPUP"`
	PageTitle   *string `json:"page_title"`
	PageURL     *string `json:"page_url"`
}

type UpdateContactStatusDTO struct {
	Status string `json:"status" validate:"required,oneof=NEW CONTACTED QUALIFIED CLOSED SPAM"`
}

type ContactResponseDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CompanyName *string   `json:"company_name"`
	Email       *string   `json:"email"`
	Phone       string    `json:"phone"`
	Country     *string   `json:"country"`
	InquiryType *string   `json:"inquiry_type"`
	ProductName *string   `json:"product_name"`
	Quantity    *string   `json:"quantity"`
	Message     *string   `json:"message"`
	Source      string    `json:"source"`
	PageTitle   *string   `json:"page_title"`
	PageURL     *string   `json:"page_url"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateContactResponseDTO struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

type GenericResponseDTO struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GetContactsQueryDTO struct {
	Page   int    `validate:"min=1"`
	Limit  int    `validate:"min=1,max=100"`
	Search string `validate:"max=255"`
	Status string `validate:"omitempty,oneof=NEW CONTACTED QUALIFIED CLOSED SPAM"`
	Source string `validate:"omitempty,oneof=CONTACT_PAGE FOOTER POPUP"`
}

type ContactPaginationDTO struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type GetAllContactsResponseDTO struct {
	Contacts   []ContactResponseDTO `json:"contacts"`
	Pagination ContactPaginationDTO `json:"pagination"`
}
