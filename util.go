package mindreminder

import (
	"fmt"
	"reflect"
	"strconv"

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
func GetPrimaryKeyValue(db *gorm.DB) string {
	var res string
	if keys := db.Statement.Schema.PrimaryFields; keys != nil {
		model := reflect.Indirect(reflect.ValueOf(db.Statement.Model))
		primary := reflect.New(model.Type()).Elem()
		res = "{"
		for _, key := range keys {
			keyValue := model.FieldByName(key.Name)
			if res != "{" {
				res += ","
			}
			res += "\"" + key.Name + "\":"
			var valueString string
			switch keyValue.Kind() {
			case reflect.Int:
				valueString = strconv.Itoa(int(keyValue.Int()))
				break
			default:
				valueString = keyValue.String()
				break
			}
			res += "\"" + valueString + "\""
			primary.FieldByName(key.Name).Set(model.FieldByName(key.Name))
		}
		res += "}"
	}
	return res
}

// //restituisce la struct con solo i campi della primary key valorizzati
// func GetPrimaryKeyValue(db *gorm.DB) interface{} {
// 	if keys := db.Statement.Schema.PrimaryFields; keys != nil {
// 		model := reflect.Indirect(reflect.ValueOf(db.Statement.Model))
// 		primary := reflect.New(model.Type()).Elem()
// 		for _, k := range keys {
// 			primary.FieldByName(k.Name).Set(model.FieldByName(k.Name))
// 		}
// 		return primary.Interface()
// 	}
// 	return 0
// }
