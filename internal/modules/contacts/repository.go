package contacts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/smart-safety-hub/backend/shared"
)

type ContactRepo struct {
	db *sqlx.DB
}

func NewContactRepo(db *sqlx.DB) *ContactRepo {
	return &ContactRepo{db: db}
}

func (r *ContactRepo) SaveContact(ctx context.Context, request CreateContactRequestDTO) (*Contact, error) {

	var contact Contact

	query := `
		INSERT INTO contacts (
			name,
			company_name,
			email,
			phone,
			country,
			inquiry_type,
			product_name,
			quantity,
			message,
			source,
			page_title,
			page_url
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12
		)
		RETURNING
    id,
    name,
    company_name,
    email,
    phone,
    country,
    inquiry_type,
    product_name,
    quantity,
    message,
    source,
    page_title,
    page_url,
    status,
    created_at,
    updated_at
	`

	err := r.db.GetContext(
		ctx,
		&contact,
		query,

		request.Name,
		request.CompanyName,
		request.Email,
		request.Phone,
		request.Country,

		request.InquiryType,
		request.ProductName,
		request.Quantity,

		request.Message,

		request.Source,

		request.PageTitle,
		request.PageURL,
	)

	if err != nil {
		return nil, shared.PostgresError(err)
	}

	return &contact, nil
}

func (r *ContactRepo) GetContactByID(ctx context.Context, contactID string) (*Contact, error) {

	var contact Contact

	query := `
		SELECT
			id,
			name,
			company_name,
			email,
			phone,
			country,
			inquiry_type,
			product_name,
			quantity,
			message,
			source,
			page_title,
			page_url,
			status,
			created_at,
			updated_at
		FROM contacts
		WHERE id = $1
	`

	if err := r.db.GetContext(ctx, &contact, query, contactID); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("contact enquiry not found")
		}

		return nil, shared.PostgresError(err)
	}

	return &contact, nil
}

func (r *ContactRepo) GetAllContacts(ctx context.Context, request GetContactsQueryDTO) ([]Contact, int, error) {

	var contacts []Contact

	args := []interface{}{}
	where := []string{"1 = 1"}

	addArg := func(value interface{}) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}

	// Search
	if request.Search != "" {
		searchPlaceholder := addArg("%" + request.Search + "%")

		where = append(
			where,
			fmt.Sprintf(
				`(
					name ILIKE %s
					OR company_name ILIKE %s
					OR email ILIKE %s
					OR phone ILIKE %s
					OR product_name ILIKE %s
					OR inquiry_type ILIKE %s
					OR message ILIKE %s
				)`,
				searchPlaceholder,
				searchPlaceholder,
				searchPlaceholder,
				searchPlaceholder,
				searchPlaceholder,
				searchPlaceholder,
				searchPlaceholder,
			),
		)
	}

	// Status
	if request.Status != "" {
		statusPlaceholder := addArg(request.Status)

		where = append(where, fmt.Sprintf("status = %s", statusPlaceholder))
	}

	// Source
	if request.Source != "" {
		sourcePlaceholder := addArg(request.Source)

		where = append(where, fmt.Sprintf("source = %s", sourcePlaceholder))
	}

	whereClause := strings.Join(
		where,
		" AND ",
	)

	// COUNT

	countQuery := fmt.Sprintf(`		SELECT COUNT(*)		FROM contacts		WHERE %s	`, whereClause)

	var total int

	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, shared.PostgresError(err)
	}

	// PAGINATION

	offset := (request.Page - 1) * request.Limit

	limitPlaceholder := addArg(request.Limit)
	offsetPlaceholder := addArg(offset)

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			company_name,
			email,
			phone,
			country,
			inquiry_type,
			product_name,
			quantity,
			message,
			source,
			page_title,
			page_url,
			status,
			created_at,
			updated_at
		FROM contacts
		WHERE %s
		ORDER BY created_at DESC
		LIMIT %s
		OFFSET %s
	`,
		whereClause,
		limitPlaceholder,
		offsetPlaceholder,
	)

	if err := r.db.SelectContext(ctx, &contacts, query, args...); err != nil {
		return nil, 0, shared.PostgresError(err)
	}

	return contacts, total, nil
}

func (r *ContactRepo) UpdateContactStatus(ctx context.Context, contactID string, status string) error {
	query := `UPDATE contacts SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2	`

	result, err := r.db.ExecContext(ctx, query, status, contactID)
	if err != nil {
		return shared.PostgresError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("contact enquiry not found")
	}

	return nil
}
