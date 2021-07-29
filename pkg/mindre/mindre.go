package mindre

import (
	"reflect"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/utils"
	"gorm.io/gorm"
)

//true se value implementa l'interfaccia Remindable e se è abilitato
func IsRemindable(value interface{}) bool {
	_, ok := value.(remindable)
	return ok
}

// Remind implementa l'interfaccia Remindable
type Remind struct {
}

func (l *Remind) AfterCreate(db *gorm.DB) error {
	return addRecordRemindToCalculate(db, utils.ActionCreate)
}

func (l *Remind) AfterUpdate(db *gorm.DB) error {
	return addRecordRemindToCalculate(db, utils.ActionUpdate)
}

func (l *Remind) AfterDelete(db *gorm.DB) error {
	return addRecordRemindToCalculate(db, utils.ActionDelete)
}

func StartService(structList []interface{}, appName string, db *gorm.DB) error {
	if _, err := config.Create(db); err != nil {
		return err
	}
	registerTypes(structList)
	if err := ricalcolaScadenze(); err != nil {
		return err
	}
	return nil
}

func NewBaseReminder(l interface{}, description string, remindType string) (Reminder, error) {
	raw, err := utils.StructToMap(&l)
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
