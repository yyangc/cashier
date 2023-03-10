package db

import (
	"context"
	"log"
	"testing"
	"time"

	"cashier/internal/model"
	"cashier/internal/model/query"
	iDB "cashier/internal/repository/database"

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

type PromotionSuite struct {
	suite.Suite

	ctx  context.Context
	repo iDB.IDatabase
}

func TestPromotion(t *testing.T) {
	suite.Run(t, new(PromotionSuite))
}

func (p *PromotionSuite) SetupSuite() {
	dsn := "default:root@tcp(127.0.0.1:3306)/default?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	p.Require().NoError(err)

	p.ctx = context.Background()
	p.repo = New(db, db)
}

func (p *PromotionSuite) TestListPromotions() {
	mps, err := p.repo.ListPromotions(p.ctx, &query.PromotionOptions{})
	p.Require().NoError(err)
	for i := range mps {
		log.Printf("%+v", mps[i])
		log.Printf("ext: %+v", mps[i].Extension)
	}
}

func (p *PromotionSuite) TestCreatePromotion() {

	mpMember := &model.Promotion{
		Name:        "pppp",
		Description: "ddddddd",
		Type:        model.PromotionTypeMember,
		Extension: &model.PromotionExtMember{
			MemberRatio: map[model.MemberType]map[int8]decimal.Decimal{
				model.MemberTypeVIP: {
					1: decimal.NewFromFloat(0.9),
					2: decimal.NewFromFloat(0.92),
					3: decimal.NewFromFloat(0.95),
				},
				model.MemberTypePro: {
					1: decimal.NewFromFloat(0.8),
					2: decimal.NewFromFloat(0.82),
					3: decimal.NewFromFloat(0.85),
				},
			},
		},
		IsDefault: true,
		StartAt:   time.Now().Add(-5 * 24 * time.Hour),
		EndAt:     time.Now().Add(15 * 24 * time.Hour),
	}

	err := p.repo.CreatePromotion(p.ctx, mpMember)
	p.Require().NoError(err)

	mpPoint := &model.Promotion{
		Name:        "pointt",
		Description: "pointtddddddd",
		Type:        model.PromotionTypePoint,
		Extension: &model.PromotionExtPoint{
			RatioToPoint: decimal.NewFromFloat(1.1),
		},
		IsDefault: true,
		StartAt:   time.Now().Add(-5 * 24 * time.Hour),
		EndAt:     time.Now().Add(15 * 24 * time.Hour),
	}

	err = p.repo.CreatePromotion(p.ctx, mpPoint)
	p.Require().NoError(err)
}
