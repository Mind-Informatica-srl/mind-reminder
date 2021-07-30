package mindre

import (
	"reflect"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/logic"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/utils"
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

func NewBaseRemind(l interface{}, description string, remindType string) (Remind, error) {
	raw, err := utils.StructToMap(&l)
	if err != nil {
		return Remind{}, err
	}
	return Remind{
		Description: &description,
		RemindType:  remindType,
		ObjectRaw:   raw,
		ObjectType:  reflect.TypeOf(l).Elem().Name(),
	}, nil
}
