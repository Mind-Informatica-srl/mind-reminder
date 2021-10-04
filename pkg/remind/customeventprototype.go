package remind

import (
	"time"

	"github.com/lib/pq"
)

// CustomEventPrototype prototipo evento di tipo custom
type CustomEventPrototype struct {
	ID int
	// nome del prototipo
	Name string
	// descrizione del prototipo
	Description *string
	// // campo per indicare il riferimento dell'evento
	// // se risorsa umana, oggetto custom o altro
	// Reference string
	// data di creazione del prototipo
	CreatedAt time.Time
	// data di modifica del prototipo
	UpdatedAt time.Time
	// jsonb con i dati dell'evento
	PrototypeEventData PrototypeEventData
	// template per la descrizione dell'evento
	EventDescriptionTemplate string
	// chiave per EventDate per generazione Event
	EventDateKey string
	// chiave per AccomplishMinScore per generazione Event
	AccomplishMinScoreKey string
	// chiave per AccomplishMaxScore per generazione Event
	AccomplishMaxScoreKey string
	// chiave per ExpectedScore per generazione Event
	ExpectedScoreKey string
	// elenco di chiavi per Hook per generazione Event
	HookKeys pq.StringArray `gorm:"type:text[]"`
	// chiave per RemindExpirationDate per generazione Event
	RemindExpirationDateKey string
	// chiave per RemindType per generazione Event
	RemindTypeKey string
	// chiave per RemindMaxScore per generazione Event
	RemindMaxScoreKey string
	// template per la descrizione della scadenza per generazione Event
	RemindDescriptionTemplate string
	// template per la descrizione dell'oggetto della scadenza per generazione Event
	RemindObjectDescriptionTemplate string
	// elenco di chiavi per RemindHook per generazione Event
	RemindHookKeys pq.StringArray `gorm:"type:text[]"`
}

// SetPK set the pk for the model
func (c *CustomEventPrototype) SetPK(pk interface{}) error {
	c.ID = pk.(int)
	return nil
}

// VerifyPK check the pk value
func (c *CustomEventPrototype) VerifyPK(pk interface{}) (bool, error) {
	id := pk.(int)
	return c.ID == id, nil
}
