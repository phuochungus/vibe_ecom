package service

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	apperrors "golf-store/be-mono/internal/shared/errors"
	"golf-store/be-mono/internal/shared/model"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

type ListInput struct {
	Query     string
	Status    string
	MinCents  *int64
	MaxCents  *int64
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
	AdminView bool
}

type ListOutput struct {
	Items      []*model.Product
	Page       int
	PageSize   int
	Total      int
	TotalPages int
}

type AdminUpsertInput struct {
	SKU         string
	Name        string
	Description string
	PriceCents  int64
	Stock       int
	Status      model.ProductStatus
	ImageURL    string
}

func (s *Service) List(input ListInput) ListOutput {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}

	where := []string{"deleted_at IS NULL"}
	args := make([]any, 0)

	if !input.AdminView {
		where = append(where, "status = 'ACTIVE'")
	}
	if strings.TrimSpace(input.Status) != "" {
		where = append(where, "status = ?")
		args = append(args, strings.ToUpper(strings.TrimSpace(input.Status)))
	}
	if strings.TrimSpace(input.Query) != "" {
		q := "%" + strings.ToLower(strings.TrimSpace(input.Query)) + "%"
		where = append(where, "(LOWER(name) LIKE ? OR LOWER(sku) LIKE ?)")
		args = append(args, q, q)
	}
	if input.MinCents != nil {
		where = append(where, "price_cents >= ?")
		args = append(args, *input.MinCents)
	}
	if input.MaxCents != nil {
		where = append(where, "price_cents <= ?")
		args = append(args, *input.MaxCents)
	}

	whereSQL := strings.Join(where, " AND ")
	sortBy := sanitizeSortBy(input.SortBy)
	sortOrder := sanitizeSortOrder(input.SortOrder)

	var total int
	countQuery := "SELECT COUNT(*) FROM products WHERE " + whereSQL
	if err := s.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return ListOutput{Items: []*model.Product{}, Page: input.Page, PageSize: input.PageSize, Total: 0, TotalPages: 0}
	}

	offset := (input.Page - 1) * input.PageSize
	listQuery := fmt.Sprintf(
		`SELECT id, sku, name, description, price_cents, stock, status, image_url, created_at, updated_at, deleted_at
		   FROM products
		  WHERE %s
		  ORDER BY %s %s
		  LIMIT ? OFFSET ?`,
		whereSQL, sortBy, sortOrder,
	)
	listArgs := append(args, input.PageSize, offset)
	rows, err := s.db.Query(listQuery, listArgs...)
	if err != nil {
		return ListOutput{Items: []*model.Product{}, Page: input.Page, PageSize: input.PageSize, Total: total, TotalPages: calcTotalPages(total, input.PageSize)}
	}
	defer rows.Close()

	items := make([]*model.Product, 0)
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			continue
		}
		items = append(items, p)
	}

	return ListOutput{
		Items:      items,
		Page:       input.Page,
		PageSize:   input.PageSize,
		Total:      total,
		TotalPages: calcTotalPages(total, input.PageSize),
	}
}

func (s *Service) GetByID(productID string, adminView bool) (*model.Product, *apperrors.APIError) {
	query := `SELECT id, sku, name, description, price_cents, stock, status, image_url, created_at, updated_at, deleted_at
	            FROM products
	           WHERE id = ? AND deleted_at IS NULL`
	args := []any{productID}
	if !adminView {
		query += ` AND status = 'ACTIVE'`
	}

	row := s.db.QueryRow(query, args...)
	product, err := scanProduct(row)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query product"}
	}

	return product, nil
}

func (s *Service) AdminCreate(input AdminUpsertInput) (*model.Product, *apperrors.APIError) {
	now := time.Now().UTC()
	if input.PriceCents <= 0 {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "price must be greater than 0"}
	}
	if input.Stock < 0 {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "stock cannot be negative"}
	}

	id := uuid.NewString()
	_, err := s.db.Exec(
		`INSERT INTO products (id, sku, name, description, price_cents, stock, status, image_url, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		strings.TrimSpace(input.SKU),
		strings.TrimSpace(input.Name),
		strings.TrimSpace(input.Description),
		input.PriceCents,
		input.Stock,
		input.Status,
		strings.TrimSpace(input.ImageURL),
		now,
		now,
	)
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "invalid product payload"}
	}

	return s.GetByID(id, true)
}

func (s *Service) AdminUpdate(productID string, input AdminUpsertInput) (*model.Product, *apperrors.APIError) {
	updates := make([]string, 0)
	args := make([]any, 0)

	if strings.TrimSpace(input.SKU) != "" {
		updates = append(updates, "sku = ?")
		args = append(args, strings.TrimSpace(input.SKU))
	}
	if strings.TrimSpace(input.Name) != "" {
		updates = append(updates, "name = ?")
		args = append(args, strings.TrimSpace(input.Name))
	}
	if strings.TrimSpace(input.Description) != "" {
		updates = append(updates, "description = ?")
		args = append(args, strings.TrimSpace(input.Description))
	}
	if input.PriceCents > 0 {
		updates = append(updates, "price_cents = ?")
		args = append(args, input.PriceCents)
	}
	if input.Stock >= 0 {
		updates = append(updates, "stock = ?")
		args = append(args, input.Stock)
	}
	if input.Status != "" {
		updates = append(updates, "status = ?")
		args = append(args, input.Status)
	}
	if strings.TrimSpace(input.ImageURL) != "" {
		updates = append(updates, "image_url = ?")
		args = append(args, strings.TrimSpace(input.ImageURL))
	}
	if len(updates) == 0 {
		return s.GetByID(productID, true)
	}

	updates = append(updates, "updated_at = ?")
	args = append(args, time.Now().UTC())
	args = append(args, productID)

	query := `UPDATE products SET ` + strings.Join(updates, ", ") + ` WHERE id = ? AND deleted_at IS NULL`
	res, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "invalid product payload"}
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return nil, apperrors.ErrNotFound
	}

	return s.GetByID(productID, true)
}

func (s *Service) AdminDelete(productID string) *apperrors.APIError {
	now := time.Now().UTC()
	res, err := s.db.Exec(
		`UPDATE products
		    SET deleted_at = ?, status = 'DISCONTINUED', updated_at = ?
		  WHERE id = ? AND deleted_at IS NULL`,
		now, now, productID,
	)
	if err != nil {
		return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to delete product"}
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func ParseStatus(status string) (model.ProductStatus, error) {
	if status == "" {
		return model.ProductStatusActive, nil
	}
	up := strings.ToUpper(strings.TrimSpace(status))
	switch up {
	case string(model.ProductStatusActive):
		return model.ProductStatusActive, nil
	case string(model.ProductStatusInactive):
		return model.ProductStatusInactive, nil
	case string(model.ProductStatusDiscontinued):
		return model.ProductStatusDiscontinued, nil
	default:
		return "", fmt.Errorf("invalid product status")
	}
}

func scanProduct(scanner interface {
	Scan(dest ...any) error
}) (*model.Product, error) {
	product := &model.Product{}
	var status string
	var deletedAt sql.NullTime

	if err := scanner.Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Description,
		&product.PriceCents,
		&product.Stock,
		&status,
		&product.ImageURL,
		&product.CreatedAt,
		&product.UpdatedAt,
		&deletedAt,
	); err != nil {
		return nil, err
	}

	product.Status = model.ProductStatus(status)
	if deletedAt.Valid {
		t := deletedAt.Time.UTC()
		product.DeletedAt = &t
	}
	return product, nil
}

func sanitizeSortBy(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "price":
		return "price_cents"
	case "created_at":
		return "created_at"
	default:
		return "name"
	}
}

func sanitizeSortOrder(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "desc") {
		return "DESC"
	}
	return "ASC"
}

func calcTotalPages(total int, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}
