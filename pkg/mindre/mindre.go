package mindre

import (
	"reflect"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/logic"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/utils"
	mrmodel "github.com/Mind-Informatica-srl/mind-reminder/pkg/mrnodel"
	"gorm.io/gorm"
)

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

func RemindToCalculateFromCreate(element interface{}, objectID string, db *gorm.DB) error {
	return logic.AddRecordRemindToCalculate(element, objectID, mrmodel.ActionCreate, db)
}

func RemindToCalculateFromUpdate(element interface{}, objectID string, db *gorm.DB) error {
	return logic.AddRecordRemindToCalculate(element, objectID, mrmodel.ActionUpdate, db)
}

func RemindToCalculateFromDelete(element interface{}, objectID string, db *gorm.DB) error {
	return logic.AddRecordRemindToCalculate(element, objectID, mrmodel.ActionDelete, db)
}
