package controllers

import (
	"errors"
	"reflect"
	"time"

	models "github.com/Mind-Informatica-srl/mind-reminder/pkg/models"
	"gorm.io/gorm"
)

//resituisce righe da lavorare (elaborated_at is null) in ordine di data creazione decrescente (in modo da avere le struct in object_raw più aggiornate)
//per ogni coppia object_id, object_type restituisce solo la riga inserita più recentemente
//si trascurano quindi le righe più vecchie (verranno comunque aggiornate dal servizio)
//(se abbiamo due righe non ancora lavorate con stesso object_id e object_type non ha infatti senso eseguire due volte il ricalcolo delle scadenze.
//Faremo il ricalcolo prendendo solo il rigo con created_at più recente)
func GetRemindToCalculate(db *gorm.DB) ([]models.RemindToCalculate, error) {
	var list []models.RemindToCalculate
	if err := db.Select("object_id, object_type").Group("object_id, object_type").Where("elaborated_at is null").Find(&list).Error; err != nil {
		return nil, err
	}
	for i, v := range list {
		if err := db.Order("created_at desc").Scopes(filterNotElaborated(v.ObjectID, v.ObjectType)).First(&v).Error; err != nil {
			return nil, err
		}
		list[i] = v
	}
	return list, nil
}

//si aggiornano tutte le righe di remind_to_calculate che hanno stesso object_id, object_type, created_at precedente a timeStart e elaborated_at null
func UpdateCorrelatedRemindToCalculate(db *gorm.DB, toCalculate *models.RemindToCalculate, timeStart time.Time, errorString *string) error {
	if err := db.Model(models.RemindToCalculate{}).Scopes(filterNotElaborated(toCalculate.ObjectID, toCalculate.ObjectType)).Where("created_at < ?", timeStart).Updates(models.RemindToCalculate{ElaboratedAt: &timeStart, Error: errorString}).Error; err != nil {
		return err
	}
	return nil
}

func filterNotElaborated(objectID string, objectType string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("object_id = ? and object_type = ? and elaborated_at is null", objectID, objectType)
	}
}

//restituisce ObjectRaw (json) sotto forma di struct sfruttando typeRegistry per ricavare il tipo di struct da ObjectType
func GetObjectFromRemindToCalculate(el *models.RemindToCalculate, typeRegistry map[string]reflect.Type) (models.Remindable, error) {
	// si ricava il tipo di struct da ObjectType
	if t, ok := typeRegistry[el.ObjectType]; ok {
		//si converte ObjectRaw (json) in struct e si mette dentro Object di el
		if err := el.PrepareObject(t); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Missing ObjectType " + el.ObjectType + " in typeRegistry")
	}
	//si esegue il cast per vedere che implementi correttamente Remindable
	obj, ok := el.Object.(models.Remindable)
	if !ok {
		return nil, errors.New("Error in cast el.Object in models.Remindable")
	}
	return obj, nil
}
