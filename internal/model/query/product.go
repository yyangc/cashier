package query

type ProductOptions struct {
	IDIn []int64

	Lock       bool
	LockNoWait bool

	// true 查詢 model.Product 關聯的 model.Inventory 並返回
	// false 則不查詢 model.Inventory
	WithInventory bool
}
