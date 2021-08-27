package logic

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/utils"
	mrmodel "github.com/Mind-Informatica-srl/mind-reminder/pkg/mrnodel"
	"github.com/Mind-Informatica-srl/restapi/pkg/models"
)

// RemindToCalculate è struct delle modifiche da cui si dovranno calcolare le scadenze (Reminder)
type RemindToCalculate struct {
	// Primary key of reminders.
	ID int //   uuid.UUID `gorm:"primary_key;"`
	// azione che ha scatenato l'insert del rigo
	Action mrmodel.Action
	// ID of tracking object.
	// By this ID later you can find all object (database row) changes.
	ObjectID string // `gorm:"index"`
	// Reflect name of tracking object.
	// It does not use package or module name, so
	// it may be not unique when use multiple types from different packages but with the same name.
	ObjectType string // `gorm:"index"`
	// Raw representation of tracking object.
	ObjectRaw models.JSONB `gorm:"type:json"`
	// Timestamp, when reminder was created.
	CreatedAt time.Time `gorm:"default:now()"`
	// data della lavorazione del calcolo della scadenza
	ElaboratedAt *time.Time
	// eventuale messaggio di errore nel calcolo della scadenza
	Error *string
	// Field Object would contain prepared structure, parsed from RawObject as json.
	// Use RegObjectType to register object types.
	Object interface{} `gorm:"-" sql:"-"`
}

// TableName return the remind to calculate table name
func (r *RemindToCalculate) TableName() string {
	return "remind_to_calculate"
}

// NewRemindToCalculate restituisce uno slice di scadenze
func NewRemindToCalculate(
	element interface{},
	objectID string,
	action mrmodel.Action,
) (r RemindToCalculate, err error) {
	if rawObjectString, err := utils.StructToMap(element); err == nil {
		r = RemindToCalculate{
			Action:     action,
			ObjectID:   objectID,
			ObjectType: reflect.TypeOf(element).Elem().Name(),
			ObjectRaw:  rawObjectString,
		}
	}
	return
}

// Event converte il json object_raw in struct e lo mette dentro Object
func (r *RemindToCalculate) Event(objType reflect.Type) (event mrmodel.Event, err error) {
	obj := reflect.New(objType).Interface()
	data, err := json.Marshal(r.ObjectRaw) // Convert to a json string

	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	event = obj.(mrmodel.Event)
	return
}
