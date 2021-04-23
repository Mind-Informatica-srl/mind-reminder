package models

import (
	mindreminder "github.com/Mind-Informatica-srl/mind-reminder"
	"gorm.io/gorm"
)

//implementa l'interfaccia Remindable
type Remind struct {
}

func (l Remind) Reminders(db *gorm.DB) (toInsert []Reminder, toDelete []Reminder, err error) {
	toInsert = []Reminder{}
	toDelete = []Reminder{}
	err = nil
	return
}

func (l *Remind) AfterCreate(db *gorm.DB) error {
	return addRecord(db, mindreminder.ActionCreate)
}

func (l *Remind) AfterUpdate(db *gorm.DB) error {
	return addRecord(db, mindreminder.ActionUpdate)
}

func (l *Remind) AfterDelete(db *gorm.DB) error {
	return addRecord(db, mindreminder.ActionDelete)
}

// Writes new reminder row to db.
func addRecord(db *gorm.DB, action string) error {
	r, err := newRemindToCalculate(db, action)
	if err != nil {
		return nil
	}
	return db.Model(&r).Create(&r).Error
}
