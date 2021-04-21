package mindreminder

import (
	"fmt"
	"reflect"
)

const remindableTag = "gorm-remindable"

//restituisce slice dei nomi dei campi che influenzano la rigenerazione di una scadenza
//sono tutti quei campi nella struct che hanno il tag gorm-remindable=true
//questo metodo viene usato SOLO se opts.computeDiff è true
func getRemindableFieldNames(value interface{}) []string {
	var names []string

	t := reflect.TypeOf(value)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value, ok := field.Tag.Lookup(remindableTag)
		if !ok || value != "true" {
			continue
		}

		names = append(names, field.Name)
	}
	return names
}

func interfaceToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprint(v)
	}
}
