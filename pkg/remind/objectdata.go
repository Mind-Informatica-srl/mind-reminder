package remind

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// PrototypeObjectData json base per PrototypeObjectData in struct eventi ed oggetti custom
type PrototypeObjectData struct {
	Fields []DataField `json:"fields"`
}

// GormDataType per PrototypeObjectData
// serve a gorm per sapere il tipo sql
func (PrototypeObjectData) GormDataType() string {
	return "jsonb"
}

// Value Make the PrototypeObjectData struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a PrototypeObjectData) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Make the PrototypeObjectData struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *PrototypeObjectData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
