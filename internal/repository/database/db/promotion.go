package db

import (
	"context"
	"time"

	"cashier/internal/model"
	"cashier/internal/model/query"
	"cashier/internal/pkg/errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Promotion 優惠活動
type promotion struct {
	ID          int64               `gorm:"column:id"`          // ID
	Name        string              `gorm:"column:name"`        // 活動名稱
	Description string              `gorm:"column:description"` // 詳情
	Type        model.PromotionType `gorm:"column:type"`        // 活動類型
	Extension   datatypes.JSON      `gorm:"column:extension"`   // 活動內容
	IsDefault   bool                `gorm:"column:is_default"`  // 是否為預設活動
	StartAt     time.Time           `gorm:"column:start_at"`    // 活動開始時間
	EndAt       time.Time           `gorm:"column:end_at"`      // 活動結束時間
	CreatedAt   time.Time           `gorm:"column:created_at"`  // 創建時間
	UpdatedAt   time.Time           `gorm:"column:updated_at"`  // 更新時間
}

func (p promotion) TableName() string {
	return "promotions"
}

func newPromotion(mPromotion *model.Promotion) (*promotion, error) {
	extB, err := mPromotion.ToExtByte()
	if err != nil {
		return nil, err
	}

	return &promotion{
		ID:          mPromotion.ID,
		Name:        mPromotion.Name,
		Description: mPromotion.Description,
		Type:        mPromotion.Type,
		Extension:   extB,
		IsDefault:   mPromotion.IsDefault,
		StartAt:     mPromotion.StartAt,
		EndAt:       mPromotion.EndAt,
	}, nil
}

func (p *promotion) ConvertToModel() (mp *model.Promotion, err error) {
	mp = &model.Promotion{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Type:        p.Type,
		IsDefault:   p.IsDefault,
		StartAt:     p.StartAt,
		EndAt:       p.EndAt,
		CreatedAt:   p.CreatedAt,
	}

	mp.Extension, err = mp.FromExtByteTo(p.Extension)
	return
}

func buildPromotionWhereCondition(db *gorm.DB, options *query.PromotionOptions) *gorm.DB {
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

	if len(options.TypeIn) > 0 {
		values := make([]interface{}, 0, len(options.TypeIn))
		for i := range options.TypeIn {
			values = append(values, options.TypeIn[i])
		}
		clauses = append(clauses, clause.IN{
			Column: "type",
			Values: values,
		})
	}

	if options.StartAtGte != nil {
		clauses = append(clauses, clause.Gte{
			Column: "start_at",
			Value:  options.StartAtGte,
		})
	}

	if options.EndAtLt != nil {
		clauses = append(clauses, clause.Gte{
			Column: "end_at",
			Value:  options.EndAtLt,
		})
	}

	db = db.Clauses(clauses...)

	return db
}

func (db *database) CreatePromotion(ctx context.Context, mPromotion *model.Promotion) error {
	_promotion, err := newPromotion(mPromotion)
	if err != nil {
		return err
	}

	if err := db.WriteDB(ctx).Create(_promotion).Error; err != nil {
		return errors.Wrapf(duplicateOrInternalError(err), "%+v", err)
	}

	return nil
}

func (db *database) ListPromotions(ctx context.Context, options *query.PromotionOptions) ([]*model.Promotion, error) {
	var _promotions = make([]*promotion, 0)
	if err := buildPromotionWhereCondition(db.ReadDB(ctx), options).Find(&_promotions).Error; err != nil {
		return nil, errors.Wrapf(errors.ErrInternalServerError, "%+v", err)
	}

	mPromotions := make([]*model.Promotion, 0, len(_promotions))
	for i := range _promotions {
		mp, err := _promotions[i].ConvertToModel()
		if err != nil {
			return nil, err
		}
		mPromotions = append(mPromotions, mp)
	}

	return mPromotions, nil
}
