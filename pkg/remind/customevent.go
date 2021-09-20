package remind

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/models"
	"gorm.io/gorm"
)

// CustomEvent evento di tipo custom
type CustomEvent struct {
	ID int
	// id del prototipo
	CustomEventPrototypeID int
	// jsonb con i dati dell'evento (semplicemnete chiave-valore)
	EventData models.JSONB
	// prototipo dell'evento
	CustomEventPrototype CustomEventPrototype `gorm:"association_autoupdate:false;"`
	CustomObjectID       *int
	// un evento può avere associato solo un oggetto
	CustomObject *CustomObject `gorm:"association_autoupdate:false;"`
}

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
	// chiave per EventType per generazione Event
	EventTypeKey string
	// chiave per EventDate per generazione Event
	EventDateKey string
	// chiave per AccomplishMinScore per generazione Event
	AccomplishMinScoreKey string
	// chiave per AccomplishMaxScore per generazione Event
	AccomplishMaxScoreKey string
	// chiave per ExpectedScore per generazione Event
	ExpectedScoreKey string
	// elenco di chiavi per Hook per generazione Event
	HookKeys []string
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
	RemindHookKeys []string
}

// GetEvent dato EventData in c (*CustomEvent)
// restituisce un evento ed un eventuale errore
func (c *CustomEvent) GetEvent() (event Event, err error) {
	var stringValue string
	var intValue int
	var dateValue time.Time
	var ok bool
	stringValue, ok = c.EventData[c.CustomEventPrototype.EventTypeKey].(string)
	if ok {
		event.EventType = stringValue
	} else {
		err = NewCustomEventError("EventType", c.CustomEventPrototype.EventTypeKey, c)
		return
	}
	dateValue, ok = c.EventData[c.CustomEventPrototype.EventDateKey].(time.Time)
	if ok {
		event.EventDate = &dateValue
	} else {
		err = NewCustomEventError("EventDate", c.CustomEventPrototype.EventDateKey, c)
		return
	}
	intValue, ok = c.EventData[c.CustomEventPrototype.AccomplishMinScoreKey].(int)
	if ok {
		event.AccomplishMinScore = intValue
	} else {
		err = NewCustomEventError("AccomplishMinScore", c.CustomEventPrototype.AccomplishMinScoreKey, c)
		return
	}
	intValue, ok = c.EventData[c.CustomEventPrototype.AccomplishMaxScoreKey].(int)
	if ok {
		event.AccomplishMaxScore = intValue
	} else {
		err = NewCustomEventError("AccomplishMaxScore", c.CustomEventPrototype.AccomplishMaxScoreKey, c)
		return
	}
	intValue, ok = c.EventData[c.CustomEventPrototype.ExpectedScoreKey].(int)
	if ok {
		event.ExpectedScore = intValue
	} else {
		err = NewCustomEventError("ExpectedScore", c.CustomEventPrototype.ExpectedScoreKey, c)
		return
	}
	event.Hook = make(map[string]interface{}, len(c.CustomEventPrototype.HookKeys))
	for i := 0; i < len(c.CustomEventPrototype.HookKeys); i++ {
		val := c.CustomEventPrototype.HookKeys[i]
		event.Hook[val] = c.EventData[val]
	}
	dateValue, ok = c.EventData[c.CustomEventPrototype.RemindExpirationDateKey].(time.Time)
	if ok {
		event.RemindInfo.ExpirationDate = &dateValue
	} else {
		err = NewCustomEventError("RemindExpirationDate", c.CustomEventPrototype.RemindExpirationDateKey, c)
		return
	}
	stringValue, ok = c.EventData[c.CustomEventPrototype.RemindTypeKey].(string)
	if ok {
		event.RemindInfo.RemindType = stringValue
	} else {
		err = NewCustomEventError("RemindType", c.CustomEventPrototype.RemindTypeKey, c)
		return
	}
	intValue, ok = c.EventData[c.CustomEventPrototype.RemindMaxScoreKey].(int)
	if ok {
		event.RemindInfo.RemindMaxScore = intValue
	} else {
		err = NewCustomEventError("RemindMaxScore", c.CustomEventPrototype.RemindMaxScoreKey, c)
		return
	}
	stringValue, err = c.parseTemplate(c.CustomEventPrototype.RemindDescriptionTemplate)
	if err != nil {
		return
	}
	event.RemindInfo.RemindDescription = stringValue
	stringValue, err = c.parseTemplate(c.CustomEventPrototype.RemindObjectDescriptionTemplate)
	if err != nil {
		return
	}
	event.RemindInfo.ObjectDescription = stringValue
	event.RemindHook = make(map[string]interface{}, len(c.CustomEventPrototype.RemindHookKeys))
	for i := 0; i < len(c.CustomEventPrototype.RemindHookKeys); i++ {
		val := c.CustomEventPrototype.RemindHookKeys[i]
		event.RemindHook[val] = c.EventData[val]
	}
	return
}

// parseTemplate parsa templateString prendendo i dati da EventData
// ES: {{.Count}} items are made of {{.Material}}
func (c *CustomEvent) parseTemplate(templateString string) (value string, err error) {
	var tmpl *template.Template
	tmpl, err = template.New("reminddescription").Parse(templateString)
	if err != nil {
		return
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, c.EventData)
	if err != nil {
		return
	}
	value = buffer.String()
	return
}

// AfterCreate di CustomEvent
func (c *CustomEvent) AfterCreate(tx *gorm.DB) (err error) {
	var event Event
	event, err = c.GetEvent()
	if err != nil {
		return
	}
	return AddEvent(tx, &event)
}

// BeforeDelete di CustomEvent
func (c *CustomEvent) BeforeDelete(tx *gorm.DB) (err error) {
	var event Event
	event, err = c.GetEvent()
	if err != nil {
		return
	}
	return DeleteEvent(tx, &event)
}

// AfterUpdate di CustomEvent
func (c *CustomEvent) AfterUpdate(tx *gorm.DB) (err error) {
	var event Event
	event, err = c.GetEvent()
	if err != nil {
		return
	}
	return UpdateEvent(tx, &event)
}

// CustomEventError has detected invalid data
type CustomEventError struct {
	Key         string       `json:"key"`
	KeyValue    string       `json:"key_value"`
	CustomEvent *CustomEvent `json:"custom_event"`
}

// NewCustomEventError restituisce nuovo NewCustomEventError passando una stringa in input
func NewCustomEventError(key string, keyValue string, c *CustomEvent) CustomEventError {
	return CustomEventError{
		Key:         key,
		KeyValue:    keyValue,
		CustomEvent: c,
	}
}

// Error serve perchè CustomEventError implementi Errors
func (a CustomEventError) Error() string {
	return fmt.Sprintf("CustomEvent error: invalid KEY VALUE %v for KEY %v", a.KeyValue, a.Key)
}
