package v1

import (
	"errors"
	"reflect"
	"time"

	mindlogger "github.com/Mind-Informatica-srl/mind-logger"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/api/v1/controllers"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/models"
	"gorm.io/gorm"
)

func RicalcolaScadenze(appLog *mindlogger.AppLogger) error {
	if err := config.Env.GetDb(appLog).Transaction(func(db *gorm.DB) error {
		//si ricava l'orario attuale (servirà per l'update su remind_to_calculate)
		timeStart := time.Now()
		//si ricavano le righe da remind_to_calculate che non sono ancora state lavorate
		toCalculateList, err := controllers.GetRemindToCalculate(db)
		if err != nil {
			return err
		}
		//si lavora ogni rigo della lista
		for _, toCalculate := range toCalculateList {
			var errorString *string
			if err := updateReminders(db, toCalculate); err != nil {
				//si scrive l'eventuale errore nella colonna error del solo rigo lavorato
				e := err.Error()
				errorString = &e
			}
			//si aggiornano tutte le righe di remind_to_calculate che hanno stesso object_id, object_type, created_at precedente a timeStart e elaborated_at null
			if err := db.Model(models.RemindToCalculate{}).Where("object_id = ? and object_type = ? and created_at < ?", toCalculate.ObjectID, toCalculate.ObjectType, timeStart).Updates(models.RemindToCalculate{ElaboratedAt: &timeStart, Error: errorString}).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func updateReminders(db *gorm.DB, el models.RemindToCalculate) error {
	//si ricava object dal rigo
	if t, ok := typeRegistry[el.ObjectType]; ok {
		if err := el.PrepareObject(t); err != nil {
			return err
		}
	} else {
		return errors.New("Missing ObjectType or typeRegistry")
	}
	//si ricavano le scadenze da cancellare e quelle da inserire
	obj, ok := el.Object.(models.Remindable)
	if !ok {
		return errors.New("Error in cast el.Object in models.Remindable")
	}
	toInsertList, toDeleteList, err := obj.Reminders(db)
	if err != nil {
		return err
	}
	db.Transaction(func(tx2 *gorm.DB) error {
		//si cancellano ed inseriscono le scadenze ottenute
		if toDeleteList != nil && len(toDeleteList) > 0 {
			if err := tx2.Delete(toDeleteList).Error; err != nil {
				return err
			}
		}
		if toInsertList != nil && len(toInsertList) > 0 {
			if err := tx2.Create(toInsertList).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

var typeRegistry = make(map[string]reflect.Type)

// typeRegistry["MyStruct"] = reflect.TypeOf(MyStruct{})
func RegisterTypes(myTypes []interface{}) {
	for _, v := range myTypes {
		typeRegistry[reflect.TypeOf(v).Name()] = reflect.TypeOf(v)
	}
}
