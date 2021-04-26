package models

import (
	"reflect"

	mindreminder "github.com/Mind-Informatica-srl/mind-reminder"
	"gorm.io/gorm"
)

//implementa l'interfaccia Remindable
type Remind struct {
}

func (l *Remind) Reminders(db *gorm.DB) (toInsert []Reminder, toDelete []Reminder, err error) {
	toInsert = []Reminder{}
	toDelete = []Reminder{}
	err = nil
	return
}

func (l *Remind) AfterCreate(db *gorm.DB) error {
	return addRecordRemindToCalculate(db, mindreminder.ActionCreate)
}

func (l *Remind) AfterUpdate(db *gorm.DB) error {
	return addRecordRemindToCalculate(db, mindreminder.ActionUpdate)
}

func (l *Remind) AfterDelete(db *gorm.DB) error {
	return addRecordRemindToCalculate(db, mindreminder.ActionDelete)
}

// Writes new reminder row to db.
func addRecordRemindToCalculate(db *gorm.DB, action string) error {
	r, err := newRemindToCalculate(db, action)
	if err != nil {
		return nil
	}
	return db.Model(&r).Create(&r).Error
}

func NewBaseReminder(l interface{}, description string, remindType string) (Reminder, error) {
	raw, err := mindreminder.InterfaceToJsonString(&l)
	if err != nil {
		return Reminder{}, err
	}
	return Reminder{
		Description:  &description,
		ReminderType: remindType,
		ObjectRaw:    raw,
		ObjectType:   reflect.TypeOf(l).Elem().Name(),
	}, nil
}
