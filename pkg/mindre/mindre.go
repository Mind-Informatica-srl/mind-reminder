package mindre

import (
	"github.com/Mind-Informatica-srl/mind-reminder/internal/calc"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/model"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/utils"
	"gorm.io/gorm"
)

type Reminder model.Reminder

//true se value implementa l'interfaccia Remindable e se è abilitato
func IsRemindable(value interface{}) bool {
	_, ok := value.(model.Remindable)
	return ok
}

// Remind implementa l'interfaccia Remindable
type Remind struct {
}

func (l *Remind) Reminders(db *gorm.DB) (toInsert []Reminder, toDelete []Reminder, err error) {
	toInsert = []Reminder{}
	toDelete = []Reminder{}
	err = nil
	return
}

func (l *Remind) AfterCreate(db *gorm.DB) error {
	return calc.AddRecordRemindToCalculate(db, utils.ActionCreate)
}

func (l *Remind) AfterUpdate(db *gorm.DB) error {
	return calc.AddRecordRemindToCalculate(db, utils.ActionUpdate)
}

func (l *Remind) AfterDelete(db *gorm.DB) error {
	return calc.AddRecordRemindToCalculate(db, utils.ActionDelete)
}

func NewBaseReminder(l interface{}, description string, remindType string) (Reminder, error) {
	e, err := calc.NewBaseReminder(l, description, remindType)
	return Reminder(e), err
}

func StartService(structList []interface{}, appName string, db *gorm.DB) error {
	if _, err := config.Create(db); err != nil {
		return err
	}
	calc.RegisterTypes(structList)
	if err := calc.RicalcolaScadenze(); err != nil {
		return err
	}
	return nil
}
