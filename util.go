package mindreminder

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

func interfaceToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprint(v)
	}
}

//restituisce la struct con solo i campi della primary key valorizzati
func GetPrimaryKeyValue(db *gorm.DB) interface{} {
	if keys := db.Statement.Schema.PrimaryFields; keys != nil {
		model := reflect.Indirect(reflect.ValueOf(db.Statement.Model))
		primary := reflect.New(model.Type()).Elem()
		for _, k := range keys {
			primary.FieldByName(k.Name).Set(model.FieldByName(k.Name))
		}
		return primary.Interface()
	}
	return 0
}
