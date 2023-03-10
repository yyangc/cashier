package db

import (
	"context"
	"log"
	"testing"

	"cashier/internal/model"
	"cashier/internal/model/query"
	"cashier/internal/model/updates"
	iDB "cashier/internal/repository/database"

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

type ProductSuite struct {
	suite.Suite

	ctx  context.Context
	repo iDB.IDatabase
}

func TestProduct(t *testing.T) {
	suite.Run(t, new(ProductSuite))
}

func (p *ProductSuite) SetupSuite() {
	dsn := "default:root@tcp(127.0.0.1:3306)/default?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	p.Require().NoError(err)

	p.ctx = context.Background()
	p.repo = New(db, db)
}

func (p *ProductSuite) TestListProducts() {
	products, err := p.repo.ListProducts(p.ctx, &query.ProductOptions{
		IDIn: []string{"ab5n8efrirohkn9ucdgh"},
	})
	p.Require().NoError(err)
	for i := range products {
		log.Printf("%+v", products[i])
	}
}

func (p *ProductSuite) TestUpdateProducts() {
	err := p.repo.UpdateProduct(p.ctx,
		&query.ProductOptions{
			IDIn: []string{"ab5n8efrirohkn9ucdgh"},
		},
		&updates.Product{
			InventoryQuantity: &model.QuantityOperation{
				Operation: model.NumericOperationSub,
				Quantity:  1,
			},
		},
	)
	p.Require().NoError(err)
}
