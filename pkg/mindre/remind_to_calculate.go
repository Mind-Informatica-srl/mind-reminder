package mindre

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/utils"
	"github.com/Mind-Informatica-srl/restapi/pkg/models"
	"gorm.io/gorm"
)

// struct delle modifiche da cui si dovranno calcolare le scadenze (Reminder)
type RemindToCalculate struct {
	// Primary key of reminders.
	ID int //uuid.UUID `gorm:"primary_key;"`
	// azione che ha scatenato l'insert del rigo
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
	ObjectRaw models.JSONB `gorm:"type:json"`
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

//restituisce uno slice di scadenze
func NewRemindToCalculate(db *gorm.DB, action string) (RemindToCalculate, error) {
	rawObjectString, err := utils.StructToMap(db.Statement.Model)
	if err != nil {
		return RemindToCalculate{}, err
	}
	return RemindToCalculate{
		Action:     action,
		ObjectID:   utils.GetPrimaryKeyValue(db),
		ObjectType: db.Statement.Schema.ModelType.Name(),
		ObjectRaw:  rawObjectString,
	}, nil
}

//converte il json object_raw in struct e lo mette dentro Object
func (r *RemindToCalculate) PrepareObject(objType reflect.Type) error {
	obj := reflect.New(objType).Interface()
	data, err := json.Marshal(r.ObjectRaw) // Convert to a json string

	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &obj)
	r.Object = obj
	return err
}
