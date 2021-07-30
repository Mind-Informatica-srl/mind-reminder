package mindre

import (
	"github.com/Mind-Informatica-srl/mind-reminder/internal/logic"
	"gorm.io/gorm"
)

type Reminder logic.Reminder

type Event logic.Event

type EventBase logic.EventBase

type A struct{}

func (e *EventBase) AfterCreate(db *gorm.DB) error {
	var le logic.EventBase = logic.EventBase(*e)
	return (&le).AfterCreate(db)
}

func (e *EventBase) AfterUpdate(db *gorm.DB) error {
	le := logic.EventBase(*e)
	return (&le).AfterUpdate(db)
}

func (e *EventBase) AfterDelete(db *gorm.DB) error {
	le := logic.EventBase(*e)
	return (&le).AfterDelete(db)
}

const (
	ActionCreate = logic.ActionCreate
	ActionUpdate = logic.ActionUpdate
	ActionDelete = logic.ActionDelete
)
