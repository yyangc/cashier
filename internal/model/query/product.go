package query

type ProductOptions struct {
	IDIn []string

	Lock       bool
	LockNoWait bool
}
