// Package mrmodel expose the model to use remind package with
package mrmodel

import "gorm.io/gorm"

// Action is the type representing what you can do with a record: create, update, delete
type Action string

const (
	// ActionCreate represent the creating action
	ActionCreate Action = "create"
	// ActionUpdate represent the updating action
	ActionUpdate Action = "update"
	// ActionDelete represent the deleting action
	ActionDelete Action = "delete"
)

// Reminder rappresenta la modifica delle scadenze da applicare
type Reminder interface {
	// ModifyReminds è la funzione invocata per ricalcolare le scandenze
	// in consegnuenza all'evento
	ModifyReminds(db *gorm.DB, action Action) (err error)
}

// Event è l'evento che scatena la necessità di ricalcolare le scadenze
type Event interface {
	// Reminders restituisce la lista di reminder da applicare
	Reminders(*gorm.DB) ([]Reminder, error)
	PrepareAddRecord(db *gorm.DB) (err error)
}
