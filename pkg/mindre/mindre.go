package mindre

import (
	"github.com/Mind-Informatica-srl/mind-reminder/internal/calc"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/model"
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
	return calc.AddRecordRemindToCalculate(db, config.ActionCreate)
}

func (l *Remind) AfterUpdate(db *gorm.DB) error {
	return calc.AddRecordRemindToCalculate(db, config.ActionUpdate)
}

func (l *Remind) AfterDelete(db *gorm.DB) error {
	return calc.AddRecordRemindToCalculate(db, config.ActionDelete)
}

func StartService(structList []interface{}, appName string) error {

	calc.RegisterTypes(structList)
	if err := calc.RicalcolaScadenze(); err != nil {
		return err
	}
	return nil
}
