package query

type InventoryOptions struct {
	ProductIDIn []int64

	Lock       bool
	LockNoWait bool
}
