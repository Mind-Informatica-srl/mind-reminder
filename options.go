package mindreminder

import "reflect"

type Option func(options *options)

type options struct {
	// lazyUpdate       bool
	// lazyUpdateFields []string
	objectTypes map[string]reflect.Type
	// computeDiff      bool
}

// RegObjectType maps object to type name, that is used in field Type of Remindable struct.
// On read change log operations, if plugin finds registered object type, by its name from db,
// it unmarshal field RawObject to Object field via json.Unmarshal.
//
// To access decoded object, e.g. `ReallyFunnyClient`, use type casting: `mindreminder.Object.(ReallyFunnyClient)`.
func RegObjectType(objectType string, objectStruct interface{}) Option {
	return func(options *options) {
		options.objectTypes[objectType] = reflect.Indirect(reflect.ValueOf(objectStruct)).Type()
	}
}
