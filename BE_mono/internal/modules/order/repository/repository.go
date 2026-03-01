package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	entities "golf-store/be-mono/internal/platform/entities"
)

type ListFilter struct {
	UserID   string
	Status   string
	From     *time.Time
	To       *time.Time
	Page     int
	PageSize int
	Admin    bool
}

type Repository interface {
	FindIdempotentOrder(userID string, idemKey string) (*entities.Order, error)
	CreateOrderTx(order *entities.Order, items []entities.OrderItem, tracking *entities.OrderTrackingEvent, notification *entities.Notification, productUpdates map[string]int) error
	List(filter ListFilter) ([]entities.Order, int64, error)
	FindByID(orderID string) (*entities.Order, error)
	CancelOrderTx(orderID string, userID string, updates map[string]any, tracking *entities.OrderTrackingEvent, notification *entities.Notification) error
	UpdateOrderStatusTx(orderID string, updates map[string]any, tracking *entities.OrderTrackingEvent, notification *entities.Notification) error
	MarkPaymentResultTx(orderID string, updates map[string]any, tracking *entities.OrderTrackingEvent, notification *entities.Notification) error
	LoadTracking(orderID string) ([]entities.OrderTrackingEvent, error)
	FindUserByID(userID string) (*entities.User, error)
	LockProduct(tx *gorm.DB, productID string) (*entities.Product, error)
	GetDbConn() *gorm.DB
}

type GormRepository struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) GetDbConn() *gorm.DB {
	return r.db
}

func (r *GormRepository) FindIdempotentOrder(userID string, idemKey string) (*entities.Order, error) {
	var existing entities.Order
	err := r.db.Where("user_id = ? AND idempotency_key = ?", userID, idemKey).Take(&existing).Error
	return &existing, err
}

func (r *GormRepository) CreateOrderTx(order *entities.Order, items []entities.OrderItem, tracking *entities.OrderTrackingEvent, notification *entities.Notification, productUpdates map[string]int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update stock
		for productID, quantity := range productUpdates {
			if err := tx.Model(&entities.Product{}).
				Where("id = ?", productID).
				Updates(map[string]any{
					"stock":      gorm.Expr("stock - ?", quantity),
					"updated_at": order.CreatedAt,
				}).Error; err != nil {
				return err
			}
		}

		// 2. Create order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 3. Create items
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		// 4. Create tracking
		if tracking != nil {
			if err := tx.Create(tracking).Error; err != nil {
				return err
			}
		}

		// 5. Create notification
		if notification != nil {
			if err := upsertNotification(tx, notification); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GormRepository) List(filter ListFilter) ([]entities.Order, int64, error) {
	query := r.db.Model(&entities.Order{})
	if !filter.Admin {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Status != "" {
		query = query.Where("order_status = ?", filter.Status)
	}
	if filter.From != nil {
		query = query.Where("placed_at >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("placed_at < ?", *filter.To)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	rows := make([]entities.Order, 0)
	if err := query.Order("placed_at DESC").Limit(filter.PageSize).Offset(offset).Find(&rows).Error; err != nil {
		return nil, total, err
	}

	return rows, total, nil
}

func (r *GormRepository) FindByID(orderID string) (*entities.Order, error) {
	var orderEntity entities.Order
	if err := r.db.Preload("Items").Where("id = ?", orderID).Take(&orderEntity).Error; err != nil {
		return nil, err
	}
	return &orderEntity, nil
}

func (r *GormRepository) CancelOrderTx(orderID string, userID string, updates map[string]any, tracking *entities.OrderTrackingEvent, notification *entities.Notification) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock order
		var orderEntity entities.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", orderID, userID).
			Take(&orderEntity).Error; err != nil {
			return err
		}

		if !canCancel(orderEntity.OrderStatus) {
			return errors.New("ORDER_CANNOT_CANCEL")
		}

		// 2. Update order
		if err := tx.Model(&entities.Order{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
			return err
		}

		// 3. Restore stock
		items := make([]entities.OrderItem, 0)
		if err := tx.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
			return err
		}
		for _, item := range items {
			if err := tx.Model(&entities.Product{}).
				Where("id = ?", item.ProductID).
				Updates(map[string]any{
					"stock":      gorm.Expr("stock + ?", item.Quantity),
					"updated_at": tracking.OccurredAt,
				}).Error; err != nil {
				return err
			}
		}

		// 4. Tracking & Notification
		if tracking != nil {
			if err := tx.Create(tracking).Error; err != nil {
				return err
			}
		}
		if notification != nil {
			if err := upsertNotification(tx, notification); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormRepository) UpdateOrderStatusTx(orderID string, updates map[string]any, tracking *entities.OrderTrackingEvent, notification *entities.Notification) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock
		var orderEntity entities.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", orderID).
			Take(&orderEntity).Error; err != nil {
			return err
		}

		toStatus, ok := updates["order_status"].(string)
		if ok && !isValidTransition(orderEntity.OrderStatus, toStatus) {
			return errors.New("INVALID_TRANSITION")
		}

		// 2. Update
		if err := tx.Model(&entities.Order{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
			return err
		}

		// 3. Track & Notify
		if tracking != nil {
			if err := tx.Create(tracking).Error; err != nil {
				return err
			}
		}
		if notification != nil {
			if err := upsertNotification(tx, notification); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormRepository) MarkPaymentResultTx(orderID string, updates map[string]any, tracking *entities.OrderTrackingEvent, notification *entities.Notification) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var orderEntity entities.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", orderID).
			Take(&orderEntity).Error; err != nil {
			return err
		}

		if err := tx.Model(&entities.Order{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
			return err
		}

		if tracking != nil {
			if err := tx.Create(tracking).Error; err != nil {
				return err
			}
		}
		if notification != nil {
			if err := upsertNotification(tx, notification); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormRepository) LoadTracking(orderID string) ([]entities.OrderTrackingEvent, error) {
	rows := make([]entities.OrderTrackingEvent, 0)
	if err := r.db.Where("order_id = ?", orderID).Order("occurred_at ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *GormRepository) FindUserByID(userID string) (*entities.User, error) {
	var user entities.User
	err := r.db.Select("id", "status").Where("id = ?", userID).Take(&user).Error
	return &user, err
}

func (r *GormRepository) LockProduct(tx *gorm.DB, productID string) (*entities.Product, error) {
	var product entities.Product
	dbConn := tx
	if dbConn == nil {
		dbConn = r.db
	}
	err := dbConn.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", productID).Take(&product).Error
	return &product, err
}

// Helper functions that we extracted from service

func canCancel(status string) bool {
	switch status {
	case entities.OrderStatusNew, entities.OrderStatusPendingPayment, entities.OrderStatusPaid, entities.OrderStatusProcessing:
		return true
	default:
		return false
	}
}

func isValidTransition(from string, to string) bool {
	if from == to {
		return true
	}
	allowed := map[string][]string{
		entities.OrderStatusNew:            {entities.OrderStatusPendingPayment, entities.OrderStatusCancelled},
		entities.OrderStatusPendingPayment: {entities.OrderStatusPaid, entities.OrderStatusCancelled},
		entities.OrderStatusPaid:           {entities.OrderStatusProcessing, entities.OrderStatusCancelled},
		entities.OrderStatusProcessing:     {entities.OrderStatusShipping, entities.OrderStatusCancelled},
		entities.OrderStatusShipping:       {entities.OrderStatusCompleted},
		entities.OrderStatusCompleted:      {},
		entities.OrderStatusCancelled:      {},
	}
	for _, candidate := range allowed[from] {
		if candidate == to {
			return true
		}
	}
	return false
}

func upsertNotification(tx *gorm.DB, notification *entities.Notification) error {
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "event_type"},
			{Name: "event_key"},
		},
		DoUpdates: clause.Assignments(map[string]any{"updated_at": notification.UpdatedAt}),
	}).Create(notification).Error
}
