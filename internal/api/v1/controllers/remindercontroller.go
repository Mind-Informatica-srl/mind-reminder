package controllers

import (
	"reflect"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/models"
	"gorm.io/gorm"
)

func UpdateReminders(db *gorm.DB, el *models.RemindToCalculate, typeRegistry map[string]reflect.Type) error {
	//si ricava object da RemindToCalculate
	obj, err := GetObjectFromRemindToCalculate(el, typeRegistry)
	if err != nil {
		return err
	}
	//si ricavano le scadenze da cancellare e quelle da inserire
	toInsertList, toDeleteList, err := obj.Reminders(db)
	if err != nil {
		return err
	}
	//si apre transazione: se una sola insert o una sola delete ha sollevato errore, si fa rollback
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
