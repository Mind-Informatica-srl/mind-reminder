package remind

import "gorm.io/gorm"

type ReferenceObject interface {
	GetID() interface{}
	GetDescription() string
}

var referenceObjectMap map[string](func(db *gorm.DB, id interface{}) (ReferenceObject, error))

// RegisterReferenceObject assegna f a objectType
// f funzione che restituisce un ReferenceObject
func RegisterReferenceObject(objectType string, f func(db *gorm.DB, id interface{}) (ReferenceObject, error)) {
	if referenceObjectMap == nil {
		referenceObjectMap = make(map[string](func(db *gorm.DB, id interface{}) (ReferenceObject, error)))
	}
	referenceObjectMap[objectType] = f
}

func GetReferenceObject(db *gorm.DB, objType string, id interface{}) (obj *ReferenceObject, err error) {
	f := referenceObjectMap[objType]
	if f != nil {
		*obj, err = f(db, id)
	}
	return
}
