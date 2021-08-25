package mindre

import (
	"github.com/Mind-Informatica-srl/mind-reminder/internal/logic"
	mrmodel "github.com/Mind-Informatica-srl/mind-reminder/pkg/mrnodel"
	"gorm.io/gorm"
)

// EventBase implementa i metodi di trigger aggiungendo un record alla
// tabella degli eventi da calcolare
type EventBase struct {
}

func (l *EventBase) PrepareAddRecord(db *gorm.DB) (err error) {
	return
}

func (l *EventBase) AfterCreate(db *gorm.DB) error {
	ev := mrmodel.Event(l)
	if err := ev.PrepareAddRecord(db); err != nil {
		return err
	}
	return logic.AddRecordRemindToCalculate(db, mrmodel.ActionCreate)
}

func (l *EventBase) AfterUpdate(db *gorm.DB) error {
	ev := mrmodel.Event(l)
	if err := ev.PrepareAddRecord(db); err != nil {
		return err
	}
	return logic.AddRecordRemindToCalculate(db, mrmodel.ActionUpdate)
}

func (l *EventBase) AfterDelete(db *gorm.DB) error {
	ev := mrmodel.Event(l)
	if err := ev.PrepareAddRecord(db); err != nil {
		return err
	}
	return logic.AddRecordRemindToCalculate(db, mrmodel.ActionDelete)
}

func (l *EventBase) Reminders(*gorm.DB) (reminders []mrmodel.Reminder, err error) {
	return
}
