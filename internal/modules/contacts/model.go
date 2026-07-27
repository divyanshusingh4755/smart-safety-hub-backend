package contacts

import "time"

type Contact struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	CompanyName *string   `db:"company_name"`
	Email       *string   `db:"email"`
	Phone       string    `db:"phone"`
	Country     *string   `db:"country"`
	InquiryType *string   `db:"inquiry_type"`
	ProductName *string   `db:"product_name"`
	Quantity    *string   `db:"quantity"`
	Message     *string   `db:"message"`
	Source      string    `db:"source"`
	PageTitle   *string   `db:"page_title"`
	PageURL     *string   `db:"page_url"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
