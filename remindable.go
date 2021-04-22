package mindreminder

import (
	"encoding/json"
	"reflect"
	"time"

	"gorm.io/gorm"
)

//struct delle scadenze
type Reminder struct {
	ID int
	//Descrizione della scadenza
	Description string
	//Tipo della scadenza
	ReminderType string
	//json dell'oggetto
	RawObject string
	//tipo della model dell'oggetto
	ObjectType string
	//Data Scadenza
	ExpireAt time.Time
	//Data Assolvenza
	AccomplishedAt *time.Time
	//Data Creazione
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	//Percentuale di assolvenza
	Percentage float64
	//Descrizione dello stato della scadenza
	StatusDescription *string
	//criteri di visibilià
	Visibility *string
}

func (t *Reminder) TableName() string {
	return "reminder"
}

// Remindable is used to get metadata from models.
type Remindable interface {
	// deve restituire slice dei reminder da inserire, slice dei reminder da cancellare e l'eventualem errore
	Reminders(*gorm.DB) ([]Reminder, []Reminder, error)
	// check if callback enabled
	isEnabled() bool
	// enable/disable loggable
	Enable(v bool)
}

type Remind struct {
	Disabled bool `gorm:"-" sql:"-" json:"-"`
}

func (l Remind) Reminders(db *gorm.DB) (toInsert []Reminder, toDelete []Reminder, err error) {
	toInsert = []Reminder{}
	toDelete = []Reminder{}
	err = nil
	return
}
func (l Remind) isEnabled() bool { return !l.Disabled }
func (l Remind) Enable(v bool)   { l.Disabled = !v }

//struct delle modifiche da cui si dovranno calcolare le scadenze (Reminder)
type ToRemind struct {
	// Primary key of reminders.
	ID int //uuid.UUID `gorm:"primary_key;"`
	// Timestamp, when reminder was created.
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
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
	// Field Object would contain prepared structure, parsed from RawObject as json.
	// Use RegObjectType to register object types.
	Object interface{} `sql:"-"`
}

func (t *ToRemind) TableName() string {
	return "to_remind"
}

//true se value implementa l'interfaccia Remindable e se è abilitato
func isRemindable(value interface{}) bool {
	v, ok := value.(Remindable)
	return ok && v.isEnabled()
}

// Allocate new and try to decode reminder field RawObject to Object.
func (l *ToRemind) prepareObject(objType reflect.Type) error {
	obj := reflect.New(objType).Interface()
	err := json.Unmarshal([]byte(l.RawObject), obj)
	l.Object = obj
	return err
}
