package db

import (
	"cashier/internal/model"
	iDB "cashier/internal/repository/database"
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
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

func (s *OrderSuite) SetupSuite() {
	readDB, writeDB, err := newTestDB()
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.repo = New(readDB, writeDB)
}

func (s *OrderSuite) TestCreateOrder() {
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
				ProductID: 1,
				Name:      "name_1",
				UnitPrice: decimal.NewFromInt32(50),
				Quantity:  1,
			},
			{
				OrderID:   orderID,
				ProductID: 2,
				Name:      "name_2",
				UnitPrice: decimal.NewFromInt32(50),
				Quantity:  1,
			},
		},
	}

	err := s.repo.CreateOrder(s.ctx, _order)
	s.Require().NoError(err)
}
