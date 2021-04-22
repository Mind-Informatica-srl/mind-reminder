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
		model := reflect.ValueOf(db.Statement.Model)
		t := reflect.TypeOf(model)
		primary := reflect.New(t)
		// var keysValue string
		for _, k := range keys {
			// v := reflect.Indirect(reflect.ValueOf(k))
			// if v.Kind() == reflect.String || v.Kind() == reflect.Int {
			primary.FieldByName(k.Name).Set(model.FieldByName(k.Name))
			// } else {
			// 	keysValue
			// }

			// var a = db.Statement.Model[k.Name]
			// println(a)
		}
		return primary.Interface()
	}
	return 0
}
