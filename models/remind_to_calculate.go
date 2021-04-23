package models

import (
	"encoding/json"
	"reflect"
	"time"

	mindreminder "github.com/Mind-Informatica-srl/mind-reminder"
	"gorm.io/gorm"
)

//struct delle modifiche da cui si dovranno calcolare le scadenze (Reminder)
type RemindToCalculate struct {
	// Primary key of reminders.
	ID int //uuid.UUID `gorm:"primary_key;"`
	//azione che ha scatenato l'insert del rigo
	Action string
	// ID of tracking object.
	// By this ID later you can find all object (database row) changes.
	ObjectID string //`gorm:"index"`
	// Reflect name of tracking object.
	// It does not use package or module name, so
	// it may be not unique when use multiple types from different packages but with the same name.
	ObjectType string //`gorm:"index"`
	// Raw representation of tracking object.
	// todo(@sas1024): Replace with []byte, to reduce allocations. Would be major version.
	ObjectRaw string `gorm:"type:json"`
	// Timestamp, when reminder was created.
	CreatedAt time.Time `gorm:"default:now()"`
	//data della lavorazione del calcolo della scadenza
	ElaboratedAt *time.Time
	//eventuale messaggio di errore nel calcolo della scadenza
	Error *string
	// Field Object would contain prepared structure, parsed from RawObject as json.
	// Use RegObjectType to register object types.
	Object interface{} `gorm:"-" sql:"-"`
}

func (t *RemindToCalculate) TableName() string {
	return "remind_to_calculate"
}

// Allocate new and try to decode reminder field RawObject to Object.
func (l *RemindToCalculate) PrepareObject(objType reflect.Type) error {
	obj := reflect.New(objType).Interface()
	err := json.Unmarshal([]byte(l.ObjectRaw), obj)
	l.Object = obj
	return err
}

//restituisce uno slice di scadenze
func newRemindToCalculate(db *gorm.DB, action string) (RemindToCalculate, error) {
	rawObject, err := json.Marshal(db.Statement.Model)
	if err != nil {
		return RemindToCalculate{}, err
	}
	return RemindToCalculate{
		Action:     action,
		ObjectID:   mindreminder.GetPrimaryKeyValue(db),
		ObjectType: db.Statement.Schema.ModelType.Name(),
		ObjectRaw:  string(rawObject),
	}, nil
}
