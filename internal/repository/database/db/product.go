package db

import (
	"context"
	"fmt"
	"time"

	"cashier/internal/model"
	"cashier/internal/model/query"
	"cashier/internal/model/updates"
	"cashier/internal/pkg/errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Product schema
type product struct {
	ID                string              `gorm:"column:id"`                 // Globally Unique ID
	Name              string              `gorm:"column:name"`               // 商品名稱
	Status            model.ProductStatus `gorm:"column:status"`             // 商品上下架狀態
	Price             decimal.Decimal     `gorm:"column:price"`              // 價格(單位：平台幣)
	Quantity          int32               `gorm:"column:quantity"`           // 總數量
	InventoryQuantity int32               `gorm:"column:inventory_quantity"` // 庫存數量
	CreatedAt         time.Time           `gorm:"column:created_at"`         // 創建時間
	UpdatedAt         time.Time           `gorm:"column:updated_at"`         // 更新時間
}

func (p product) TableName() string {
	return "products"
}

func (p *product) ConvertToModel() *model.Product {
	return &model.Product{
		ID:                p.ID,
		Name:              p.Name,
		Status:            p.Status,
		Price:             p.Price,
		Quantity:          p.Quantity,
		InventoryQuantity: p.InventoryQuantity,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}

type productUpdates struct {
	Name              *string              `gorm:"column:name"`               // 商品名稱
	Status            *model.ProductStatus `gorm:"column:status"`             // 商品上下架狀態
	Price             *decimal.Decimal     `gorm:"column:price"`              // 價格(單位：平台幣)
	Quantity          *gormExpr            `gorm:"column:quantity"`           // 總數量
	InventoryQuantity *gormExpr            `gorm:"column:inventory_quantity"` // 庫存數量
	UpdatedAt         *time.Time           `gorm:"column:updated_at"`         // 更新時間
}

func buildProductWhereCondition(db *gorm.DB, options *query.ProductOptions) *gorm.DB {
	var clauses []clause.Expression

	if len(options.IDIn) > 0 {
		values := make([]interface{}, 0, len(options.IDIn))
		for i := range options.IDIn {
			values = append(values, options.IDIn[i])
		}
		clauses = append(clauses, clause.IN{
			Column: "id",
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

func (db *database) ListProducts(ctx context.Context, options *query.ProductOptions) ([]*model.Product, error) {
	var products = make([]*product, 0)

	if err := buildProductWhereCondition(db.ReadDB(ctx), options).Find(&products).Error; err != nil {
		return nil, notFoundOrInternalError(err)
	}

	var mProducts = make([]*model.Product, 0, len(products))
	for i := range products {
		mProducts = append(mProducts, products[i].ConvertToModel())
	}

	return mProducts, nil
}

func (db *database) UpdateProduct(ctx context.Context, options *query.ProductOptions, updates *updates.Product) error {
	var _updates = &productUpdates{}

	if updates.InventoryQuantity != nil {
		_updates.InventoryQuantity = &gormExpr{clause.Expr{
			SQL:  fmt.Sprintf("%s %s ?", "inventory_quantity", updates.InventoryQuantity.Operation.Sql()),
			Vars: []interface{}{updates.InventoryQuantity.Quantity},
		}}
	}

	if updates.Quantity != nil {
		_updates.Quantity = &gormExpr{clause.Expr{
			SQL:  fmt.Sprintf("%s %s ?", "quantity", updates.Quantity.Operation.Sql()),
			Vars: []interface{}{updates.InventoryQuantity.Quantity},
		}}
	}

	if err := buildProductWhereCondition(db.WriteDB(ctx), options).
		Table(product{}.TableName()).
		Updates(_updates).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
