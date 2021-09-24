package remind

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// PrototypeEventData json base per PrototypeEventData in struct eventi ed oggetti custom
type PrototypeEventData struct {
	Fields []DataField `json:"fields"`
}

// GormDataType per PrototypeEventData
// serve a gorm per sapere il tipo sql
func (PrototypeEventData) GormDataType() string {
	return "jsonb"
}

// Value Make the PrototypeEventData struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a PrototypeEventData) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Make the PrototypeEventData struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *PrototypeEventData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// DataField struct per DataField usato in PrototypeEventData e PrototypeObjectData
type DataField struct {
	Name         string        `json:"name"`
	Type         DataFieldType `json:"type"`
	Label        string        `json:"label"`
	DefaultValue *interface{}  `json:"default_value"`
	Required     bool          `json:"required"`
	Tooltip      *string       `json:"tooltip"`
	Hint         *string       `json:"hint"`
	Options      *[]string     `json:"options"`
}
