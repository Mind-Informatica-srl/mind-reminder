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
