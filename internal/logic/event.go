package logic

import "gorm.io/gorm"

const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// Reminder rappresenta la modifica delle scadenze da applicare
type Reminder interface {
	// ModifyReminds è la funzione invocata per ricalcolare le scandenze
	// in consegnuenza all'evento
	ModifyReminds(db *gorm.DB, action string) error
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

// EventBase implementa i metodi di trigger aggiungendo un record alla
// tabella degli eventi da calcolare
type EventBase struct {
}

func (l *EventBase) AfterCreate(db *gorm.DB) error {
	return AddRecordRemindToCalculate(db, ActionCreate)
}

func (l *EventBase) AfterUpdate(db *gorm.DB) error {
	return AddRecordRemindToCalculate(db, ActionUpdate)
}

func (l *EventBase) AfterDelete(db *gorm.DB) error {
	return AddRecordRemindToCalculate(db, ActionDelete)
}
