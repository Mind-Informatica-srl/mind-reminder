package model

import "gorm.io/gorm"

type Remindable interface {
	// deve restituire slice dei reminder da inserire, slice dei reminder da cancellare e l'eventualem errore
	Reminders(*gorm.DB, string) ([]Reminder, []Reminder, string, error)
	//viene chiamato da gorm dopo insert
	AfterCreate(*gorm.DB) error
	//viene chiamato da gorm dopo update
	AfterUpdate(*gorm.DB) error
	//viene chiamato da gorm dopo delete
	AfterDelete(*gorm.DB) error
}
