package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/smart-safety-hub/backend/shared"
)

type ProductRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (r *ProductRepo) SaveProduct(ctx context.Context, request ProductRequestDTO) (*string, error) {
	query := "INSERT INTO products(name, slug, description, seller_id, brand_id, category_id, status) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id"
	var lastInsertId string
	if err := r.db.QueryRowContext(ctx, query, request.Name, request.Slug, request.Description, request.SellerID, request.BrandID, request.CategoryID, request.Status).Scan(&lastInsertId); err != nil {
		return nil, shared.PostgresError(err)
	}
	return &lastInsertId, nil
}

func (r *ProductRepo) UpdateProduct(ctx context.Context, productID string, request ProductRequestDTO) error {
	query := "UPDATE products SET name=COALESCE(NULLIF($1, ''), name), slug=COALESCE(NULLIF($2, ''), slug), description=COALESCE(NULLIF($3, ''), description), seller_id=COALESCE(NULLIF($4, '')::UUID, seller_id), brand_id=COALESCE(NULLIF($5, '')::UUID, brand_id), category_id=COALESCE(NULLIF($6, '')::UUID, category_id), status=COALESCE(NULLIF($7, '')::status_enum, status) WHERE id=$8"
	if _, err := r.db.ExecContext(ctx, query, request.Name, request.Slug, request.Description, request.SellerID, request.BrandID, request.CategoryID, request.Status, productID); err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *ProductRepo) DeleteProduct(ctx context.Context, productID string, status string) error {
	query := "UPDATE products SET status=$1 WHERE id=$2"
	if _, err := r.db.ExecContext(ctx, query, status, productID); err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *ProductRepo) GetProductByID(ctx context.Context, productID string) (*Product, error) {
	var product Product
	query := "SELECT * FROM products WHERE id=$1"
	if err := r.db.GetContext(ctx, &product, query, productID); err != nil {
		fmt.Println("err", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, shared.PostgresError(err)
	}
	return &product, nil
}

func (r *ProductRepo) GetProductBySlug(ctx context.Context, slug string) (*Product, error) {
	var product Product
	query := `SELECT * FROM products WHERE slug = $1 AND status='ACTIVE' LIMIT 1`

	err := r.db.GetContext(ctx, &product, query, slug)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepo) GetAllProducts(ctx context.Context, request ProductFilters) ([]GetProducts, error) {
	query := `SELECT 
		p.id, p.name, p.slug, p.description, p.status, 
		b.name AS brand_name, 
		c.name AS category_name,
		media.url AS image_url,
		COUNT(*) OVER() AS total_count
		FROM products p 
		LEFT JOIN brands b ON p.brand_id = b.id 
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN LATERAL (
		SELECT url
		FROM product_media pm
		WHERE pm.product_id = p.id AND type = 'image'
		ORDER BY display_order ASC
		LIMIT 1
		) media ON true
		 WHERE 1=1
		`

	var args []interface{}

	if request.Status != "" {
		query += " AND p.status = ?"
		args = append(args, request.Status)
	}

	if len(request.Category) > 0 {
		q, inArgs, err := sqlx.In(" AND c.slug IN (?)", request.Category)
		if err != nil {
			return nil, err
		}

		query += q
		args = append(args, inArgs...)
	}

	if len(request.Brand) > 0 {
		q, inArgs, err := sqlx.In(" AND b.slug IN (?)", request.Brand)
		if err != nil {
			return nil, err
		}

		query += q
		args = append(args, inArgs...)
	}

	if request.Search != "" {
		query += " AND (p.name ILIKE ? OR p.description ILIKE ?)"
		searchTerm := "%" + request.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	query += " ORDER BY p.created_at DESC LIMIT ? OFFSET ?"

	if request.Limit <= 0 {
		request.Limit = 40
	}
	offset := (request.Page - 1) * request.Limit
	args = append(args, request.Limit, offset)

	query = r.db.Rebind(query)

	var products []GetProducts
	if err := r.db.SelectContext(ctx, &products, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, shared.PostgresError(err)
	}

	return products, nil
}

func (r *ProductRepo) AddProductAttribute(ctx context.Context, productID string, request []ProductAttributeArrayDTO) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return shared.PostgresError((err))
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM products_attributes WHERE product_id=$1", productID)
	if err != nil {
		return shared.PostgresError(err)
	}

	if len(request) > 0 {
		numFields := 3
		placeholders := make([]string, len(request))
		values := make([]interface{}, 0, len(request)*numFields)

		for i, attr := range request {
			offset := i * numFields
			placeholders[i] = fmt.Sprintf("($%d, $%d, $%d)", offset+1, offset+2, offset+3)

			// Append values
			values = append(values, attr.ProductID, attr.AttributeKey, attr.AttributeValue)
		}

		query := fmt.Sprintf(`INSERT INTO products_attributes (product_id, attribute_key, attribute_value) VALUES %s ON CONFLICT (product_id, attribute_key) DO UPDATE SET attribute_value = EXCLUDED.attribute_value, updated_at = CURRENT_TIMESTAMP`, strings.Join(placeholders, ","))
		_, err := tx.ExecContext(ctx, query, values...)
		if err != nil {
			return shared.PostgresError(err)
		}
	}
	return tx.Commit()
}

func (r *ProductRepo) GetProductAttributeByID(ctx context.Context, productID string) ([]ProductAttribute, error) {
	var productAttribute []ProductAttribute
	query := "SELECT * FROM products_attributes WHERE product_id=$1"
	if err := r.db.SelectContext(ctx, &productAttribute, query, productID); err != nil {
		fmt.Println("err", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, shared.PostgresError(err)
	}
	return productAttribute, nil
}

func (r *ProductRepo) SyncProductVariants(ctx context.Context, productId string, req VariantRequestDTO) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return shared.PostgresError(err)
	}

	defer tx.Rollback()

	// Delete on old data
	if _, err := tx.ExecContext(ctx, "DELETE FROM product_options WHERE product_id = $1", productId); err != nil {
		return shared.PostgresError(err)
	}

	// We keep existing variants to preserve their IDs, but clear their old option associations
	if _, err := tx.ExecContext(ctx, "DELETE FROM variant_option_values WHERE variant_id IN (SELECT id FROM product_variants WHERE product_id = $1)", productId); err != nil {
		return shared.PostgresError(err)
	}

	// Insert OPTIONs and VALUES (Capturing IDs for mapping)
	valueMap := make(map[string]string)
	for _, opt := range req.Options {
		var optionID string
		if err := tx.QueryRowContext(ctx, "INSERT INTO product_options (product_id, name) VALUES ($1,$2) RETURNING id", productId, opt.Name).Scan(&optionID); err != nil {
			return shared.PostgresError(err)
		}

		for _, val := range opt.Values {
			var valID string
			if err := tx.QueryRowContext(ctx, "INSERT INTO product_option_values (option_id, value) VALUES ($1, $2) RETURNING id", optionID, val).Scan(&valID); err != nil {
				return shared.PostgresError(err)
			}
			valueMap[val] = valID
		}
	}

	// Bulk Upsert Variants
	if len(req.Variants) > 0 {
		variantPlaceholders := []string{}
		variantValues := []interface{}{}
		numFields := 5

		for i, v := range req.Variants {
			offset := i * numFields
			variantPlaceholders = append(variantPlaceholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", offset+1, offset+2, offset+3, offset+4, offset+5))
			variantValues = append(variantValues, productId, v.SKU, v.Price, v.Weight, v.IsActive)
		}
		variantQuery := fmt.Sprintf(`INSERT INTO product_variants (product_id, sku, price, weight, is_active) VALUES %s ON CONFLICT (sku) DO UPDATE SET price = EXCLUDED.price, weight = EXCLUDED.weight, is_active = EXCLUDED.is_active, updated_at = CURRENT_TIMESTAMP RETURNING id, sku`, strings.Join(variantPlaceholders, ","))
		rows, err := tx.QueryContext(ctx, variantQuery, variantValues...)
		if err != nil {
			return shared.PostgresError(err)
		}

		defer rows.Close()

		// Bulk insert linking table (variant_option_values)
		linkPlaceholders := []string{}
		linkValues := []interface{}{}
		linkCount := 1

		// Map sku to returned id for linking
		skuToID := make(map[string]string)
		for rows.Next() {
			var id, sku string
			if err := rows.Scan(&id, &sku); err != nil {
				return shared.PostgresError(err)
			}
			skuToID[sku] = id
		}

		for _, v := range req.Variants {
			vID := skuToID[v.SKU]
			for _, optValName := range v.OptionValues {
				if valID, ok := valueMap[optValName]; ok {
					linkPlaceholders = append(linkPlaceholders, fmt.Sprintf("($%d, $%d)", linkCount, linkCount+1))
					linkValues = append(linkValues, vID, valID)
					linkCount += 2
				}
			}
		}

		if len(linkPlaceholders) > 0 {
			linkQuery := fmt.Sprintf("INSERT INTO variant_option_values (variant_id, option_value_id) VALUES %s ON CONFLICT DO NOTHING", strings.Join(linkPlaceholders, ","))
			if _, err := tx.ExecContext(ctx, linkQuery, linkValues...); err != nil {
				return shared.PostgresError(err)
			}
		}
	}

	return tx.Commit()
}

func (r *ProductRepo) GetProductVariants(ctx context.Context, productId string) (*ProductVariants, error) {
	query := `
	WITH product_options_data AS (
	SELECT po.id, po.name, jsonb_agg(pov.value) AS values FROM product_options po JOIN product_option_values pov ON po.id = pov.option_id WHERE po.product_id = $1 GROUP BY po.id, po.name
	),
	variants_data AS (
	SELECT pv.id, pv.sku, pv.price, pv.weight, pv.is_active, jsonb_agg(pov.value) AS option_values FROM product_variants pv JOIN variant_option_values vov ON pv.id = vov.variant_id JOIN product_option_values pov ON vov.option_value_id = pov.id WHERE pv.product_id = $1 GROUP BY pv.id, pv.sku, pv.price, pv.weight, pv.is_active
	)
	SELECT 
	$1 AS product_id, 
	(SELECT jsonb_agg(jsonb_build_object('name', name, 'values', values )) FROM product_options_data) AS options,
	(SELECT jsonb_agg(jsonb_build_object('id', id, 'sku', sku, 'price', price, 'weight', weight, 'is_active', is_active, 'option_values', option_values)) FROM variants_data) AS variants;
	`

	var result ProductVariants

	if err := r.db.GetContext(ctx, &result, query, productId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *ProductRepo) AddProductMedia(ctx context.Context, productId string, media []ProductMediaDTO) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return shared.PostgresError(err)
	}

	defer tx.Rollback()

	// Delete old media to sync with current one
	_, err = tx.ExecContext(ctx, "DELETE FROM product_media WHERE product_id = $1", productId)
	if err != nil {
		return shared.PostgresError(err)
	}

	if len(media) == 0 {
		return tx.Commit()
	}

	// Prepare Bulk Insert
	numFields := 5
	placeholders := make([]string, len(media))
	values := make([]interface{}, 0, len(media)*numFields)

	for i, m := range media {
		offset := i * numFields
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", offset+1, offset+2, offset+3, offset+4, offset+5)

		values = append(values, productId, m.VariantID, m.Url, m.MediaType, m.DisplayOrder)
	}

	query := fmt.Sprintf(`INSERT INTO product_media (product_id, variant_id, url, type, display_order) VALUES %s`, strings.Join(placeholders, ","))
	if _, err = tx.ExecContext(ctx, query, values...); err != nil {
		return shared.PostgresError(err)
	}

	return tx.Commit()
}

func (r *ProductRepo) GetProductMedia(ctx context.Context, productId string) ([]ProductMedia, error) {
	var productMedia []ProductMedia

	query := `
		SELECT id, product_id, variant_id, url, type, display_order 
		FROM product_media 
		WHERE product_id = $1 
		ORDER BY display_order ASC`

	if err := r.db.SelectContext(ctx, &productMedia, query, productId); err != nil {
		return nil, shared.PostgresError(err)
	}

	return productMedia, nil
}

func (r *ProductRepo) SaveProductSEO(ctx context.Context, seo ProductSEO) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return shared.PostgresError(err)
	}

	defer tx.Rollback()

	query := `
		INSERT INTO product_seo (product_id, meta_title, meta_description, og_image_url, keywords)
		VALUES (:product_id, :meta_title, :meta_description, :og_image_url, :keywords)
		ON CONFLICT (product_id) DO UPDATE SET
			meta_title = EXCLUDED.meta_title,
			meta_description = EXCLUDED.meta_description,
			og_image_url = EXCLUDED.og_image_url,
			keywords = EXCLUDED.keywords`

	_, err = tx.NamedExecContext(ctx, query, seo)
	if err != nil {
		return shared.PostgresError(err)
	}

	publishProduct := `UPDATE products SET status = 'ACTIVE' WHERE id = $1`
	_, err = tx.ExecContext(ctx, publishProduct, seo.ProductID)
	if err != nil {
		return shared.PostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *ProductRepo) GetProductSEO(ctx context.Context, productId string) (*ProductSEO, error) {
	var seo ProductSEO
	query := `
		SELECT product_id, meta_title, meta_description, og_image_url, keywords 
		FROM product_seo 
		WHERE product_id = $1`

	err := r.db.GetContext(ctx, &seo, query, productId)
	if err != nil {
		return nil, shared.PostgresError(err)
	}

	return &seo, nil
}
