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

func (l *EventBase) AfterCreate(db *gorm.DB) error {
	return logic.AddRecordRemindToCalculate(db, mrmodel.ActionCreate)
}

func (l *EventBase) AfterUpdate(db *gorm.DB) error {
	return logic.AddRecordRemindToCalculate(db, mrmodel.ActionUpdate)
}

func (l *EventBase) AfterDelete(db *gorm.DB) error {
	return logic.AddRecordRemindToCalculate(db, mrmodel.ActionDelete)
}
