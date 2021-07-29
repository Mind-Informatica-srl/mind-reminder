package mindre

import (
	"reflect"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/calc"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/model"
	"gorm.io/gorm"
)

//true se value implementa l'interfaccia Remindable e se è abilitato
func IsRemindable(value interface{}) bool {
	_, ok := value.(model.Remindable)
	return ok
}

//implementa l'interfaccia Remindable
type Remind struct {
}

func (l *Remind) Reminders(db *gorm.DB) (toInsert []model.Reminder, toDelete []model.Reminder, err error) {
	toInsert = []model.Reminder{}
	toDelete = []model.Reminder{}
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

func NewBaseReminder(l interface{}, description string, remindType string) (model.Reminder, error) {
	raw, err := config.InterfaceToJsonString(&l)
	if err != nil {
		return model.Reminder{}, err
	}
	return model.Reminder{
		Description:  &description,
		ReminderType: remindType,
		ObjectRaw:    raw,
		ObjectType:   reflect.TypeOf(l).Elem().Name(),
	}, nil
}

func StartService(structList []interface{}, appName string) error {

	calc.RegisterTypes(structList)
	if err := calc.RicalcolaScadenze(); err != nil {
		return err
	}
	return nil
}
