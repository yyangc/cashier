package model

// NumericOperation 操作
type NumericOperation int8

const (
	NumericOperationUnknown NumericOperation = iota
	NumericOperationAdd                      // 加
	NumericOperationSub                      // 減
)

func (no NumericOperation) Sql() string {
	switch no {
	case NumericOperationAdd:
		return "+"
	case NumericOperationSub:
		return "-"
	}
	return ""
}

// DiscountType 折扣類型
type DiscountType int8

const (
	DiscountTypeUnknown DiscountType = iota
	// DiscountTypeRate 百分比
	DiscountTypeRate
	// DiscountTypeAmount 金額
	DiscountTypeAmount
)
