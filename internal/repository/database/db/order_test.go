package db

import (
	"context"
	"testing"

	"cashier/internal/model"
	iDB "cashier/internal/repository/database"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ################################
//
//  超級隨便的測試
//  只是想測試 sql 語法正常
//
// ################################

type OrderSuite struct {
	suite.Suite

	ctx  context.Context
	repo iDB.IDatabase
}

func TestOrder(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}

func (o *OrderSuite) SetupSuite() {
	dsn := "default:root@tcp(127.0.0.1:3306)/default?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	o.Require().NoError(err)

	o.ctx = context.Background()
	o.repo = New(db, db)
}

func (o *OrderSuite) TestCreateOrder() {
	orderID := xid.New().String()
	_order := &model.Order{
		ID:            xid.New().String(),
		UserID:        12345678,
		OriginalPrice: decimal.NewFromInt32(100),
		FinalPrice:    decimal.NewFromInt32(85),
		UsedPoints:    5,
		PromotionIDs:  []int64{111, 222},
		Items: []*model.OrderItem{
			{
				OrderID:   orderID,
				ProductID: "ProductID_1",
				Name:      "name_1",
				UnitPrice: decimal.NewFromInt32(50),
				Quantity:  1,
			},
			{
				OrderID:   orderID,
				ProductID: "ProductID_2",
				Name:      "name_2",
				UnitPrice: decimal.NewFromInt32(50),
				Quantity:  1,
			},
		},
	}

	err := o.repo.CreateOrder(o.ctx, _order)
	o.Require().NoError(err)
}
