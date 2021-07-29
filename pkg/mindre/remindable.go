package mindre

import "gorm.io/gorm"

type remindable interface {
	// deve restituire slice dei reminder da inserire, slice dei reminder da cancellare e l'eventualem errore
	Reminders(db *gorm.DB, action string) (toInsert []Reminder, toDelete []Reminder, err error)
	//viene chiamato da gorm dopo insert
	AfterCreate(*gorm.DB) error
	//viene chiamato da gorm dopo update
	AfterUpdate(*gorm.DB) error
	//viene chiamato da gorm dopo delete
	AfterDelete(*gorm.DB) error
}
