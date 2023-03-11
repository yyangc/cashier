package db

import (
	"context"
	"testing"

	"cashier/internal/model"
	"cashier/internal/model/query"
	"cashier/internal/model/updates"
	iDB "cashier/internal/repository/database"

	"github.com/stretchr/testify/suite"
)

// ################################
//
//  超級隨便的測試
//  只是想測試 sql 語法正常
//
// ################################

type InventorySuite struct {
	suite.Suite

	ctx  context.Context
	repo iDB.IDatabase
}

func TestInventory(t *testing.T) {
	suite.Run(t, new(InventorySuite))
}

func (s *InventorySuite) SetupSuite() {
	readDB, writeDB, err := newTestDB()
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.repo = New(readDB, writeDB)

}

func (s *InventorySuite) TestUpdateProducts() {
	err := s.repo.UpdateInventory(s.ctx,
		&query.InventoryOptions{
			ProductIDIn: []int64{1},
		},
		&updates.Inventory{
			AvailableQuantity: &model.QuantityOperation{
				Operation: model.NumericOperationSub,
				Quantity:  1,
			},
		},
	)
	s.Require().NoError(err)
}
