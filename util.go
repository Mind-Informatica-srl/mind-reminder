package mindreminder

import (
	"fmt"
	"reflect"
)

const remindableTag = "gorm-remindable"

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
