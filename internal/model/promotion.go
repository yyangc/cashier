package model

import (
	"cashier/internal/pkg/errors"
	"encoding/json"
	"sort"
	"time"

	"gorm.io/datatypes"

	"github.com/shopspring/decimal"
)

// ValidPromotionTypes 優惠活動類型
var ValidPromotionTypes = []PromotionType{PromotionTypeMember, PromotionTypePoint, PromotionTypeExtraDiscount}

// PromotionType 優惠類型
type PromotionType int8

const (
	PromotionTypeUnknown PromotionType = iota
	// PromotionTypeMember 會員優惠
	PromotionTypeMember
	// PromotionTypePoint 平台點數
	PromotionTypePoint
	// PromotionTypeExtraDiscount 額外優惠
	PromotionTypeExtraDiscount
)

// Promotion 優惠活動
type Promotion struct {
	ID          int64         // ID
	Name        string        // 活動名稱
	Description string        // 詳情
	Type        PromotionType // 活動類型
	Extension   IPromotionExt // 活動內容
	IsDefault   bool          // 是否為預設活動
	StartAt     time.Time     // 活動開始時間
	EndAt       time.Time     // 活動結束時間
	CreatedAt   time.Time     // 創建時間
	UpdatedAt   time.Time     // 更新時間
}

func (p *Promotion) ToExtByte() (datatypes.JSON, error) {
	b, err := json.Marshal(p.Extension)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternalError, err.Error())
	}

	return b, nil
}

func (p *Promotion) FromExtByteTo(jsonB datatypes.JSON) (IPromotionExt, error) {
	var ext IPromotionExt

	switch p.Type {
	case PromotionTypeMember:
		ext = &PromotionExtMember{}
	case PromotionTypePoint:
		ext = &PromotionExtPoint{}
	case PromotionTypeExtraDiscount:
		ext = &PromotionExtExtraDiscount{}
	}

	if err := json.Unmarshal(jsonB, ext); err != nil {
		return nil, errors.Wrap(errors.ErrInternalError, err.Error())
	}
	return ext, nil
}

type CalculatePriceInput struct {
	Member     *Member
	UsedPoints int32
}

type IPromotionExt interface {
	CalculatePrice(beforePrice decimal.Decimal, input *CalculatePriceInput) (usePromotion bool, afterPrize decimal.Decimal)
}

// PromotionExtMember 優惠類型(會員)的內容
type PromotionExtMember struct {
	// [VIP] 1 -> 95 折
	// [VIP] 2 -> 9 折
	// [VIP] 3 -> 85 折
	// [Pro] 1 -> 8 折
	MemberRatio map[MemberType]map[int8]decimal.Decimal
}

// CalculatePrice 計算優惠類型(會員)後的價格
func (p *PromotionExtMember) CalculatePrice(beforePrice decimal.Decimal, input *CalculatePriceInput) (UsedPromotion bool, afterPrice decimal.Decimal) {
	afterPrice = beforePrice
	if input.Member == nil {
		return
	}

	// 檢查是否有該 VIP 類型
	memberTypeRatio, exist := p.MemberRatio[input.Member.Type]
	if !exist {
		return
	}

	// 檢查是否有支持該 VIP 等級
	memberRatio, exist := memberTypeRatio[input.Member.Level]
	if !exist {
		return
	}

	return true, afterPrice.Mul(memberRatio)
}

// PromotionExtPoint 優惠類型(點數)的內容
type PromotionExtPoint struct {
	Ratio decimal.Decimal // 比例，平台點數:平台幣
}

// CalculatePrice 計算優惠類型(會員)後的價格
func (p *PromotionExtPoint) CalculatePrice(beforePrice decimal.Decimal, input *CalculatePriceInput) (usePromotion bool, afterPrice decimal.Decimal) {
	if input.UsedPoints == 0 {
		return false, beforePrice
	}

	return true, beforePrice.Sub(p.Ratio.Mul(decimal.NewFromInt32(input.UsedPoints)))
}

// PromotionExtExtraDiscount 優惠類型(額外優惠)的內容
type PromotionExtExtraDiscount struct {
	Requirement    ExtraDiscountRequirement // 符合條件
	DiscountType   DiscountType             // 折扣類型，e.g. 百分比、金額
	DiscountRate   decimal.Decimal          // 折抵百分比
	DiscountAmount decimal.Decimal          // 折抵金額
}

// ExtraDiscountRequirement 優惠類型(額外優惠)的符合條件
type ExtraDiscountRequirement struct {
	// [vip] 1
	// [Pro] 1, 2, 3
	MemberLevel map[MemberType][]int8
	Point       int32
}

// CalculatePrice 計算優惠類型(額外優惠)後的價格
func (p *PromotionExtExtraDiscount) CalculatePrice(beforePrice decimal.Decimal, input *CalculatePriceInput) (usePromotion bool, afterPrice decimal.Decimal) {
	afterPrice = beforePrice
	// 	額外優惠， 如果有要求會員等級
	if p.Requirement.MemberLevel != nil {
		// 用戶不是會員
		if input.Member == nil {
			return
		}

		// 用戶不符合特定的 memberType
		levels, exist := p.Requirement.MemberLevel[input.Member.Type]
		if !exist {
			return
		}

		idx := sort.Search(len(levels), func(i int) bool {
			return levels[i] == input.Member.Level
		})
		// 用戶不符合特定的等級
		if idx >= len(levels) {
			return
		}
	}

	// 	額外優惠， 如果有要求點數
	if p.Requirement.Point > 0 {
		if input.UsedPoints == 0 {
			return
		}

		// 沒有使用點數
		if input.UsedPoints == 0 {
			return
		}
	}

	// 計算額外優惠後的價格
	if p.DiscountType == DiscountTypeRate {
		afterPrice = afterPrice.Mul(p.DiscountRate)
	} else {
		afterPrice = afterPrice.Sub(p.DiscountAmount)
	}

	return true, afterPrice
}
