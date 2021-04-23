package controllers

import (
	"github.com/Mind-Informatica-srl/mind-reminder/internal/models"
	"gorm.io/gorm"
)

//resituisce righe da lavorare in ordine di data creazione decrescente (in modo da avere le struct in object_raw più aggiornate)
//per ogni coppia object_id, object_type restituisce solo la riga inserita più recentemente
func GetRemindToCalculate(db *gorm.DB) ([]models.RemindToCalculate, error) {
	var list []models.RemindToCalculate
	if err := db.Select("object_id, object_type").Group("object_id, object_type").Where("elaborated_at is null").Find(&list).Error; err != nil {
		return nil, err
	}
	for i, v := range list {
		if err := db.Order("created_at desc").Where("object_id = ? and object_type = ? and elaborated_at is null", v.ObjectID, v.ObjectType).First(&v).Error; err != nil {
			return nil, err
		}
		list[i] = v
	}
	return list, nil
}
