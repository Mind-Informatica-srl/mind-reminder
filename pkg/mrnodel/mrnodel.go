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
	ModifyReminds(db *gorm.DB, action Action) ([]Reminder, error)
}

// Event è l'evento che scatena la necessità di ricalcolare le scadenze
type Event interface {
	// Reminders restituisce la lista di reminder da applicare
	Reminders(*gorm.DB) ([]Reminder, error)
	// AfterCreate viene chiamato da gorm dopo insert
	AfterCreate(*gorm.DB) error
	// AfterUpdate viene chiamato da gorm dopo update
	AfterUpdate(*gorm.DB) error
	//AfterDelete viene chiamato da gorm dopo delete
	AfterDelete(*gorm.DB) error
}
