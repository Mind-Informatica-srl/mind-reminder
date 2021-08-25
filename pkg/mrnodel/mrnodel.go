package mrmodel

import "gorm.io/gorm"

type Action string

const (
	ActionCreate Action = "create"
	ActionUpdate        = "update"
	ActionDelete        = "delete"
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
