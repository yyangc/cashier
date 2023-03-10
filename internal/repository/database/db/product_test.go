package db

import (
	"context"
	"log"
	"testing"

	"cashier/internal/model/query"
	iDB "cashier/internal/repository/database"

	"github.com/stretchr/testify/suite"
)

// ################################
//
//  超級隨便的測試
//  只是想測試 sql 語法正常
//
// ################################

type ProductSuite struct {
	suite.Suite

	ctx  context.Context
	repo iDB.IDatabase
}

func TestProduct(t *testing.T) {
	suite.Run(t, new(ProductSuite))
}

func (s *ProductSuite) SetupSuite() {
	readDB, writeDB, err := newTestDB()
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.repo = New(readDB, writeDB)
}

func (s *ProductSuite) TestListProducts() {
	products, err := s.repo.ListProducts(s.ctx, &query.ProductOptions{
		IDIn:          []int64{},
		WithInventory: true,
	})
	s.Require().NoError(err)
	for i := range products {
		log.Printf("product: %+v", products[i])
		log.Printf("Inventory: %+v", products[i].Inventory)
	}
}
