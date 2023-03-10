package service

import (
	"cashier/internal/model"
	"cashier/internal/model/query"
	"context"
	"time"
)

func (s *service) ListPromotions(ctx context.Context, options query.PromotionOptions) ([]*model.Promotion, error) {
	promotions, err := s.db.ListPromotions(ctx, &options)
	if err != nil {
		return nil, err
	}
	return promotions, nil
}

func (s *service) GetCurrPromotionsMap(ctx context.Context) (map[model.PromotionType]*model.Promotion, error) {
	now := time.Now().UTC()
	promotions, err := s.db.ListPromotions(ctx, &query.PromotionOptions{
		TypeIn:     []model.PromotionType{model.PromotionTypeMember, model.PromotionTypePoint, model.PromotionTypeExtraDiscount},
		StartAtGte: &now,
		EndAtLt:    &now,
	})
	if err != nil {
		return nil, err
	}

	defaultPromotions := make(map[model.PromotionType]*model.Promotion, 0)    // 預設活動
	processingPromotions := make(map[model.PromotionType]*model.Promotion, 0) // 進行中的活動
	for _, p := range promotions {
		if p.IsDefault {
			defaultPromotions[p.Type] = p
		} else {
			processingPromotions[p.Type] = p
		}
	}

	var res = make(map[model.PromotionType]*model.Promotion, len(model.ValidPromotionTypes))
	for _, validType := range model.ValidPromotionTypes {
		// 如果有進行中的則優先
		if _, exist := processingPromotions[validType]; exist {
			res[validType] = processingPromotions[validType]
			continue
		}

		if _, exist := defaultPromotions[validType]; exist {
			res[validType] = defaultPromotions[validType]
			continue
		}
	}

	return res, nil
}
