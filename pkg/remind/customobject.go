package remind

import (
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/models"
)

// CustomObject oggetto di tipo custom
type CustomObject struct {
	ID int
	// id del prototipo
	CustomObjectPrototypeID int
	// jsonb con i dati dell'oggetto
	EventData models.JSONB
	// prototipo dell'oggetto
	CustomObjectPrototype *CustomObjectPrototype `gorm:"association_autoupdate:false;"`
	// un oggetto può avere associati più eventi
	CustomEvents *[]CustomEvent `gorm:"association_autoupdate:false;"`
}

// CustomObjectPrototype prototipo di oggetto di tipo custom
type CustomObjectPrototype struct {
	ID int
	// nome del prototipo
	Name string
	// descrizione del prototipo
	Description *string
	// data di creazione del prototipo
	CreatedAt time.Time
	// data di modifica del prototipo
	UpdatedAt time.Time
	// jsonb con i dati dell'evento
	PrototypeObjectData PrototypeObjectData
}
