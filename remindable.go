package mindreminder

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
	// "github.com/gofrs/uuid"
)

// Interface is used to get metadata from models.
type Interface interface {
	// Meta should return structure, that can be converted to json.
	Meta() interface{}
	// GenerateReminder should return *time.Time
	GenerateReminder() (*time.Time, error)
	// lock makes available only embedding structures.
	lock()
	// check if callback enabled
	isEnabled() bool
	// enable/disable loggable
	Enable(v bool)
}

type RemindableModel struct {
	Disabled bool `sql:"-" json:"-"`
}

func (RemindableModel) GenerateReminder() (*time.Time, error) { return nil, nil }
func (RemindableModel) Meta() interface{}                     { return nil }
func (RemindableModel) lock()                                 {}
func (l RemindableModel) isEnabled() bool                     { return !l.Disabled }
func (l RemindableModel) Enable(v bool)                       { l.Disabled = !v }

// Reminder is a main entity, which used to log changes.
// Commonly, Reminder is stored in 'reminder' table.
type Reminder struct {
	// Primary key of reminders.
	ID int //uuid.UUID `gorm:"primary_key;"`
	//timestamp indicating the time of the reminder
	RemindAt time.Time
	// Timestamp, when reminder was created.
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	// Action type.
	// On write, supports only 'create', 'update', 'delete',
	Action string
	// ID of tracking object.
	// By this ID later you can find all object (database row) changes.
	ObjectID string `gorm:"index"`
	// Reflect name of tracking object.
	// It does not use package or module name, so
	// it may be not unique when use multiple types from different packages but with the same name.
	ObjectType string `gorm:"index"`
	// Raw representation of tracking object.
	// todo(@sas1024): Replace with []byte, to reduce allocations. Would be major version.
	RawObject string `sql:"type:JSON"`
	// Raw representation of tracking object's meta.
	// todo(@sas1024): Replace with []byte, to reduce allocations. Would be major version.
	RawMeta string `sql:"type:JSON"`
	// Field Object would contain prepared structure, parsed from RawObject as json.
	// Use RegObjectType to register object types.
	Object interface{} `sql:"-"`
	// Field Meta would contain prepared structure, parsed from RawMeta as json.
	// Use RegMetaType to register object's meta types.
	Meta interface{} `sql:"-"`
}

func (t *Reminder) TableName() string {
	return "reminder"
}

func isLoggable(value interface{}) bool {
	_, ok := value.(Interface)
	return ok
}

func isEnabled(value interface{}) bool {
	v, ok := value.(Interface)
	return ok && v.isEnabled()
}

func generateReminder(scope *gorm.Scope) (*time.Time, error) {
	val, ok := scope.Value.(Interface)
	if !ok {
		return nil, errors.New("error in generateReminder")
	}
	return val.GenerateReminder()
}

func fetchChangeLogMeta(scope *gorm.Scope) []byte {
	val, ok := scope.Value.(Interface)
	if !ok {
		return nil
	}
	data, err := json.Marshal(val.Meta())
	if err != nil {
		panic(err)
	}
	return data
}

func (l *Reminder) prepareObject(objType reflect.Type) error {
	// Allocate new and try to decode reminder field RawObject to Object.
	obj := reflect.New(objType).Interface()
	err := json.Unmarshal([]byte(l.RawObject), obj)
	l.Object = obj
	return err
}

func (l *Reminder) prepareMeta(objType reflect.Type) error {
	// Allocate new and try to decode reminder field RawObject to Object.
	obj := reflect.New(objType).Interface()
	err := json.Unmarshal([]byte(l.RawMeta), obj)
	l.Meta = obj
	return err
}
