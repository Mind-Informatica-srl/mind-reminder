package utils

import (
	"encoding/json"
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

func InterfaceToJsonString(l interface{}) (string, error) {
	rawObject, err := json.Marshal(l)
	if err != nil {
		return "", err
	}
	return string(rawObject), nil
}

// Converts a struct to a map while maintaining the json alias as keys
func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}
