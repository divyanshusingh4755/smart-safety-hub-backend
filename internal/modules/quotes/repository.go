package quotes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/smart-safety-hub/backend/shared"
)

type QuoteRepo struct {
	db *sqlx.DB
}

func NewQuoteRepo(db *sqlx.DB) *QuoteRepo {
	return &QuoteRepo{
		db: db,
	}
}

func (r *QuoteRepo) SaveQuote(ctx context.Context, request CreateQuoteRequestDTO, requiredBy time.Time) (*Quote, []QuoteItem, error) {

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, shared.PostgresError(err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var quote Quote

	quoteQuery := `
		INSERT INTO quote_requests (
			quote_number,
			full_name,
			company_name,
			email,
			phone,
			delivery_location,
			required_by,
			additional_requirements
		)
		VALUES (
			'RFQ-' ||
			TO_CHAR(CURRENT_DATE, 'YYYY') ||
			'-' ||
			LPAD(
				nextval('quote_request_number_seq'::regclass)::text,
				6,
				'0'
			),
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		)
		RETURNING
			id,
			quote_number,
			full_name,
			company_name,
			email,
			phone,
			delivery_location,
			required_by,
			additional_requirements,
			status,
			created_at,
			updated_at
	`

	err = tx.GetContext(
		ctx,
		&quote,
		quoteQuery,
		request.FullName,
		request.CompanyName,
		request.Email,
		request.Phone,
		request.DeliveryLocation,
		requiredBy,
		request.AdditionalRequirements,
	)
	if err != nil {
		return nil, nil, shared.PostgresError(err)
	}

	items := make([]QuoteItem, 0, len(request.Items))

	itemQuery := `
		INSERT INTO quote_request_items (
			quote_request_id,
			product_ref,
			product_name,
			sku,
			image_url,
			quantity
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
		RETURNING
			id,
			quote_request_id,
			product_ref,
			product_name,
			sku,
			image_url,
			quantity,
			created_at
	`

	for _, item := range request.Items {
		var savedItem QuoteItem

		err = tx.GetContext(
			ctx,
			&savedItem,
			itemQuery,
			quote.ID,
			item.ProductRef,
			item.Name,
			item.SKU,
			item.Image,
			item.Quantity,
		)
		if err != nil {
			return nil, nil, shared.PostgresError(err)
		}

		items = append(items, savedItem)
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, shared.PostgresError(err)
	}

	return &quote, items, nil
}

func (r *QuoteRepo) GetQuoteByID(ctx context.Context, quoteID string) (*Quote, []QuoteItem, error) {
	var quote Quote
	quoteQuery := `SELECT id, quote_number, full_name, company_name, email, phone, delivery_location, required_by, additional_requirements, status, created_at, updated_at FROM quote_requests WHERE id = $1`

	if err := r.db.GetContext(ctx, &quote, quoteQuery, quoteID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, errors.New("quote request not found")
		}

		return nil, nil, shared.PostgresError(err)
	}

	var items []QuoteItem

	itemsQuery := `SELECT id, quote_request_id, product_ref, product_name, sku, image_url, quantity, created_at FROM quote_request_items WHERE quote_request_id = $1 ORDER BY created_at ASC`
	if err := r.db.SelectContext(ctx, &items, itemsQuery, quote.ID); err != nil {
		return nil, nil, shared.PostgresError(err)
	}

	return &quote, items, nil
}

func (r *QuoteRepo) UpdateQuoteStatus(ctx context.Context, quoteID string, status string) error {
	query := `UPDATE quote_requests SET status = $1 updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, quoteID)
	if err != nil {
		return shared.PostgresError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("quote request not found")
	}

	return nil
}

func (r *QuoteRepo) GetAllQuotes(ctx context.Context, request GetQuotesQueryDTO) ([]Quote, int, error) {
	var quotes []Quote
	args := []interface{}{}

	where := []string{"1 = 1"}
	addArg := func(value interface{}) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}

	if request.Search != "" {
		searchPlaceholder := addArg("%" + request.Search + "%")

		where = append(where,
			fmt.Sprintf(
				`(
				quote_number ILIKE %s OR full_name ILIKE %s OR company_name ILIKE %s OR email ILIKE %s OR phone ILIKE %s OR delivery_location ILIKE %s
			)`, searchPlaceholder, searchPlaceholder, searchPlaceholder, searchPlaceholder, searchPlaceholder, searchPlaceholder,
			))
	}

	if request.Status != "" {
		statusPlaceholder := addArg(request.Status)
		where = append(where, fmt.Sprintf("status = %s", statusPlaceholder))
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM quote_requests WHERE %s`, whereClause)
	var total int

	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, shared.PostgresError(err)
	}

	offset := (request.Page - 1) * request.Limit
	limitPlaceholder := addArg(request.Limit)
	offsetPlaceholder := addArg(offset)

	query := fmt.Sprintf(`
		SELECT id, quote_number, full_name, company_name, email, phone, delivery_location, required_by, additional_requirements, status, created_at, updated_at FROM quote_requests WHERE %s ORDER BY created_at DESC LIMIT %s OFFSET %s
	`, whereClause, limitPlaceholder, offsetPlaceholder)

	if err := r.db.SelectContext(ctx, &quotes, query, args...); err != nil {
		return nil, 0, shared.PostgresError(err)
	}

	return quotes, total, nil
}
