package tcc

import (
	"github.com/mushroomsir/tcc/store"
)

// Transaction ...
type Transaction struct {
	uid   string
	name  string
	store store.TccAsyncTaskInterface
}

// TryPlan get try execute play, example sql.
func (a *Transaction) TryPlan(value string) string {
	return a.store.TryPlan(a.uid, a.name, value)
}

// Try ...
func (a *Transaction) Try(value string) error {
	return a.store.Try(a.uid, a.name, value)
}

// TryAndConfirm ...
func (a *Transaction) TryAndConfirm(value string) error {
	return a.store.TryAndConfirm(a.uid, a.name, value)
}

// Confirm submit async task
func (a *Transaction) Confirm() error {
	return a.store.Confirm(a.uid)
}

// Cancel delete async task
func (a *Transaction) Cancel() error {
	return a.store.Cancel(a.uid)
}
