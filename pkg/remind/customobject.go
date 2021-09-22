package remind

import (
	"github.com/Mind-Informatica-srl/restapi/pkg/models"
)

// CustomObject oggetto di tipo custom
type CustomObject struct {
	ID int
	// id del prototipo
	CustomObjectPrototypeID int
	// jsonb con i dati dell'oggetto
	Data models.JSONB
	// prototipo dell'oggetto
	CustomObjectPrototype *CustomObjectPrototype `gorm:"association_autoupdate:false;"`
	// un oggetto può avere associati più eventi
	CustomEvents *[]CustomEvent `gorm:"association_autoupdate:false;"`
}

// SetPK set the pk for the model
func (c *CustomObject) SetPK(pk interface{}) error {
	c.ID = pk.(int)
	return nil
}

// VerifyPK check the pk value
func (c *CustomObject) VerifyPK(pk interface{}) (bool, error) {
	id := pk.(int)
	return c.ID == id, nil
}
