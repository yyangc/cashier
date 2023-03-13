package db

import (
	"context"
	"fmt"
	"time"

	"cashier/internal/model"
	"cashier/internal/model/query"
	"cashier/internal/model/updates"
	"cashier/internal/pkg/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// inventory schema
type inventory struct {
	ID                int64     `gorm:"column:id"`
	ProductID         int64     `gorm:"column:product_id"`
	TotalQuantity     int32     `gorm:"column:total_quantity"`
	AvailableQuantity int32     `gorm:"column:available_quantity"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

func (i inventory) TableName() string {
	return "inventories"
}

func (i *inventory) ConvertToModel() *model.Inventory {
	return &model.Inventory{
		ID:                i.ID,
		ProductID:         i.ProductID,
		TotalQuantity:     i.TotalQuantity,
		AvailableQuantity: i.AvailableQuantity,
		CreatedAt:         i.CreatedAt,
		UpdatedAt:         i.UpdatedAt,
	}
}

type inventoryUpdates struct {
	TotalQuantity     *gormExpr `gorm:"column:total_quantity"`
	AvailableQuantity *gormExpr `gorm:"column:available_quantity"`
}

func buildInventoryWhereCondition(db *gorm.DB, options *query.InventoryOptions) *gorm.DB {
	var clauses []clause.Expression

	if len(options.ProductIDIn) > 0 {
		values := make([]interface{}, 0, len(options.ProductIDIn))
		for i := range options.ProductIDIn {
			values = append(values, options.ProductIDIn[i])
		}
		clauses = append(clauses, clause.IN{
			Column: "product_id",
			Values: values,
		})
	}

	if options.Lock {
		var lockingOption string
		if options.LockNoWait {
			lockingOption = "NOWAIT"
		}
		clauses = append(clauses, clause.Locking{Strength: "UPDATE", Options: lockingOption})
	}

	db = db.Clauses(clauses...)

	return db
}

func (db *database) ListInventories(ctx context.Context, options *query.InventoryOptions) ([]*model.Inventory, error) {
	var _inventories = make([]*inventory, 0)

	if err := buildInventoryWhereCondition(db.ReadDB(ctx), options).Find(&_inventories).Error; err != nil {
		return nil, errors.Wrapf(duplicateOrInternalError(err), "%+v", err)
	}

	var mis = make([]*model.Inventory, 0, len(_inventories))
	for i := range _inventories {
		mis = append(mis, _inventories[i].ConvertToModel())
	}

	return mis, nil
}

func (db *database) UpdateInventory(ctx context.Context, options *query.InventoryOptions, updates *updates.Inventory) error {
	var _updates = &inventoryUpdates{}

	if updates.TotalQuantity != nil {
		_updates.TotalQuantity = &gormExpr{clause.Expr{
			SQL:  fmt.Sprintf("%s %s ?", "total_quantity", updates.TotalQuantity.Operation.Sql()),
			Vars: []interface{}{updates.TotalQuantity.Quantity},
		}}
	}

	if updates.AvailableQuantity != nil {
		_updates.AvailableQuantity = &gormExpr{clause.Expr{
			SQL:  fmt.Sprintf("%s %s ?", "available_quantity", updates.AvailableQuantity.Operation.Sql()),
			Vars: []interface{}{updates.AvailableQuantity.Quantity},
		}}
	}

	if err := buildInventoryWhereCondition(db.WriteDB(ctx), options).
		Table(inventory{}.TableName()).
		Updates(_updates).Error; err != nil {
		return errors.Wrapf(errors.ErrInternalServerError, "%+v", err)
	}

	return nil
}
