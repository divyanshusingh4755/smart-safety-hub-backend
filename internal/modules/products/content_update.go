package products

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type UpdateContentBySlugRequest struct {
	Slug            string                  `json:"slug" validate:"required"`
	Name            string                  `json:"name" validate:"required"`
	Description     string                  `json:"description" validate:"required"`
	Attributes      []ProductAttributeArray `json:"attributes" validate:"required,min=1,dive"`
	MetaTitle       string                  `json:"meta_title" validate:"required,max=255"`
	MetaDescription string                  `json:"meta_description" validate:"required"`
	Keywords        []string                `json:"keywords" validate:"required,min=1,dive,required"`
}

type UpdateContentBySlugResponse struct {
	ProductID string   `json:"product_id"`
	Slug      string   `json:"slug"`
	Status    string   `json:"status"`
	Updated   []string `json:"updated"`
}

var ErrContentProductNotFound = errors.New("product not found")

// UpdateContentBySlug performs all database writes automatically
func (r *ProductRepo) UpdateContentBySlug(ctx context.Context, req UpdateContentBySlugRequest) (string, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin transaction: %w", err)
	}

	defer tx.Rollback()

	// find and lock the product row. Unlike the public lookup, this also supports DRAFT/ARCHIVED.
	var productID string
	err = tx.GetContext(
		ctx,
		&productID,
		`SELECT id FROM products WHERE slug = $1 FOR UPDATE`,
		req.Slug,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrContentProductNotFound
	}

	if err != nil {
		return "", fmt.Errorf("lookup product: %w", err)
	}

	// Update only name and description. All other product fields remain untouched
	result, err := tx.ExecContext(
		ctx,
		`UPDATE products SET
		name = $1, description = $2, updated = CURRENT_TIMESTAMP 
		WHERE id = $3`,
		req.Name,
		req.Description,
		productID,
	)
	if err != nil {
		return "", fmt.Errorf("update prouct: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("check product update: %w", err)
	}

	if affected != 1 {
		return "", fmt.Errorf("unexpected updated row count: %d", affected)
	}

	// Replace existing attributes in the same transaction.
	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM products_attributes WHERE product_id = $1`,
		productID,
	)

	if err != nil {
		return "", fmt.Errorf("delete old attributes: %w", err)
	}

	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO products_attributes
		(product_id, attribute_key, attribute_value)
		VALUES ($1, $2, $3)`,
	)

	if err != nil {
		return "", fmt.Errorf("prepare attributes: %w", err)
	}

	defer stmt.Close()

	for _, attr := range req.Attributes {
		_, err := stmt.ExecContext(
			ctx,
			productID,
			attr.AttributeKey,
			attr.AttributeValue,
		)
		if err != nil {
			return "", fmt.Errorf("insert attribute %q: %w", attr.AttributeKey, err)
		}
	}

	// Update SEO, preserving the existing og_image_url.
	keywordsJSON, err := json.Marshal(req.Keywords)
	if err != nil {
		return "", fmt.Errorf("encode keywords: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO product_seo
		 (product_id, meta_title, meta_description, keywords)
		 VALUES ($1, $2, $3, $4::jsonb)
		 ON CONFLICT (product_id) DO UPDATE SET
		   meta_title = EXCLUDED.meta_title,
		   meta_description = EXCLUDED.meta_description,
		   keywords = EXCLUDED.keywords,
		   updated_at = CURRENT_TIMESTAMP`,
		productID,
		req.MetaTitle,
		req.MetaDescription,
		string(keywordsJSON),
	)

	if err != nil {
		return "", fmt.Errorf("update SEO: %w", err)
	}

	// All three wrties must succeed before committing.
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit content update: %w", err)
	}

	return productID, nil
}

// Service layer: invalidate caches only after commit.
func (s *ProductService) UpdateContentBySlug(ctx context.Context, req UpdateContentBySlugRequest) (*UpdateContentBySlugResponse, error) {
	productID, err := s.repo.UpdateContentBySlug(ctx, req)
	if err != nil {
		return nil, err
	}

	s.inValidateProductCache(ctx, productID, req.Slug)
	s.invalidateListCache(ctx)

	return &UpdateContentBySlugResponse{
		ProductID: productID,
		Slug:      req.Slug,
		Status:    "success",
		Updated:   []string{"product", "attributes", "seo"},
	}, nil
}

// HTTP handler for the single request content update
func (h *RestHandler) UpdateContentBySlug(w http.ResponseWriter, r *http.Request) {
	// Limit JSON request body to 2 MiB.
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req UpdateContentBySlugRequest
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Reject multiple JSON objects in one request.
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "request must contain exactly one JSON object", http.StatusBadRequest)
		return
	}

	req.Slug = strings.TrimSpace(req.Slug)
	req.Name = strings.TrimSpace(req.Name)
	req.MetaTitle = strings.TrimSpace(req.MetaTitle)
	req.MetaDescription = strings.TrimSpace(req.MetaDescription)

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Reject empty or duplicate attribute keys.
	seen := make(map[string]bool, len(req.Attributes))
	for _, attr := range req.Attributes {
		key := strings.TrimSpace(attr.AttributeKey)
		value := strings.TrimSpace(attr.AttributeValue)

		if key == "" || value == "" {
			http.Error(w, "attribute keys and values cannot be empty", http.StatusBadRequest)
			return
		}

		if seen[key] {
			http.Error(w, "duplicate attribute key: "+key, http.StatusBadRequest)
			return
		}
		seen[key] = true
	}

	response, err := h.service.UpdateContentBySlug(r.Context(), req)
	if errors.Is(err, ErrContentProductNotFound) {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		h.service.logger.Error("failed to update product content", zap.Error((err)))
		http.Error(w, "failed to update product content", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
