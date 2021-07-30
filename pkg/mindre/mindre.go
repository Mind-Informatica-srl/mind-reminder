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

// deve restituire slice dei reminder da inserire, slice dei reminder da cancellare e l'eventualem errore
func (l *Remind) Reminders(db *gorm.DB, action string) (toInsert []Reminder, toDelete []Reminder, err error) {
	// se action == insert
	// 	1. cerco le scadenze a cui l'evento assolve cronologicamente compatibili non ancora assolte o assolte da eventi successivi all'evento stesso.
	//	2. se assolte da eventi successivi all'evento stesso, cancello le assolvenze e tengo in una lista gli eventi successivi.
	//	3. inserisco l'assolvenza
	//	4. genero l'eventuale scadenza scaturita dall'evento a seconda che l'evento assolva o meno quali scadenze
	//	5. per ogni evento successivo nella lista ripeto da punto 1 (al punto 4 valuto se la scadenza già esistente corrisponde con quella che dovrei generare)

	// se action == delete
	//	1. cancello le assolvenze generate dall'evento
	//	2. metto in una lista gli eventi che assolvono alla scadenze generate dall'evento
	//	3. cancello le scadenze generate dall'evento
	//	4. per ogni evento nella lista cerco scadenze che assolve e procedo come da punto insert - 1

	// se action == update
	//	1. eseguo delete
	//	2. eseguo insert
	return
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
