// Package mindre expose the functions to configure and start the remind service
package mindre

import (
	"reflect"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/logic"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/utils"
	mrmodel "github.com/Mind-Informatica-srl/mind-reminder/pkg/mrnodel"
	"gorm.io/gorm"
)

// StartService start the remind service
func StartService(structList []interface{}, appName string, db *gorm.DB) error {
	if _, err := config.Create(db); err != nil {
		return err
	}
	logic.RegisterTypes(structList)
	if err := logic.RicalcolaScadenze(); err != nil {
		return err
	}
	return nil
}

// NewBaseRemind create a new remind with the brovided data
func NewBaseRemind(l interface{}, description string, remindType string) (mrmodel.Remind, error) {
	raw, err := utils.StructToMap(&l)
	if err != nil {
		return mrmodel.Remind{}, err
	}
	return mrmodel.Remind{
		Description: &description,
		RemindType:  remindType,
		ObjectRaw:   raw,
		ObjectType:  reflect.TypeOf(l).Elem().Name(),
	}, nil
}

// RemindToCalculateFromCreate add a new remind to calculate from a creating action
func RemindToCalculateFromCreate(element interface{}, objectID string, db *gorm.DB) error {
	return logic.AddRecordRemindToCalculate(element, objectID, mrmodel.ActionCreate, db)
}

// RemindToCalculateFromUpdate add a new remind to calculate from a updating action
func RemindToCalculateFromUpdate(element interface{}, objectID string, db *gorm.DB) error {
	return logic.AddRecordRemindToCalculate(element, objectID, mrmodel.ActionUpdate, db)
}

// RemindToCalculateFromDelete add a new remind to calculate from a deleting action
func RemindToCalculateFromDelete(element interface{}, objectID string, db *gorm.DB) error {
	return logic.AddRecordRemindToCalculate(element, objectID, mrmodel.ActionDelete, db)
}
