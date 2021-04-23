package mindreminder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
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
	var sb strings.Builder
	if keys := db.Statement.Schema.PrimaryFields; keys != nil {
		model := reflect.Indirect(reflect.ValueOf(db.Statement.Model))
		primary := reflect.New(model.Type()).Elem()
		sb.WriteString("{")
		addComma := false
		for _, key := range keys {
			keyValue := model.FieldByName(key.Name)
			if !addComma {
				addComma = true
			} else {
				sb.WriteString(",")
			}
			sb.WriteString("\"")
			sb.WriteString(key.Name)
			sb.WriteString("\":")
			var valueString string
			switch keyValue.Kind() {
			case reflect.Int:
				valueString = strconv.Itoa(int(keyValue.Int()))
				break
			default:
				valueString = keyValue.String()
				break
			}
			sb.WriteString("\"")
			sb.WriteString(valueString)
			sb.WriteString("\"")
			primary.FieldByName(key.Name).Set(model.FieldByName(key.Name))
		}
		sb.WriteString("}")
	}
	return sb.String()
}
