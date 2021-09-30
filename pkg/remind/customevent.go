package remind

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/constants"
	"github.com/Mind-Informatica-srl/restapi/pkg/models"
	"gorm.io/gorm"
)

// CustomEvent evento di tipo custom
type CustomEvent struct {
	ID int
	// id del prototipo
	CustomEventPrototypeID int
	CustomSectionID        int
	// jsonb con i dati dell'evento (semplicemnete chiave-valore)
	Data models.JSONB
	// prototipo dell'evento
	CustomEventPrototype CustomEventPrototype `gorm:"association_autoupdate:false;"`
	CustomObjectID       *int
	// un evento può avere associato solo un oggetto
	CustomObject *CustomObject `gorm:"association_autoupdate:false;"`
}

// SetPK set the pk for the model
func (c *CustomEvent) SetPK(pk interface{}) error {
	c.ID = pk.(int)
	return nil
}

// VerifyPK check the pk value
func (c *CustomEvent) VerifyPK(pk interface{}) (bool, error) {
	id := pk.(int)
	return c.ID == id, nil
}

// GetEvent dato EventData in c (*CustomEvent)
// restituisce un evento ed un eventuale errore
func (c *CustomEvent) GetEvent(db *gorm.DB) (event Event, err error) {
	var stringValue string
	var intValue int
	var dateValue time.Time
	var ok bool
	if c.CustomEventPrototype.ID == 0 {
		// se non c'è già, carico CustomEventPrototype dal db
		if err = db.Table("custom_event_prototypes").Where("id=?", c.CustomEventPrototypeID).First(&c.CustomEventPrototype).Error; err != nil {
			return
		}
	}
	if c.CustomEventPrototype.EventTypeKey == "" {
		event.EventType = strconv.Itoa(c.ID)
	} else {
		stringValue, ok = c.Data[c.CustomEventPrototype.EventTypeKey].(string)
		if ok {
			event.EventType = stringValue
		} else {
			err = NewCustomEventError("EventType", c.CustomEventPrototype.EventTypeKey, c)
			return
		}
	}
	stringValue, ok = c.Data[c.CustomEventPrototype.EventDateKey].(string)
	if ok {
		dateValue, err = ParseDate(stringValue)
		if err == nil {
			event.EventDate = &dateValue
		} else {
			err = NewCustomEventError("EventDate", c.CustomEventPrototype.EventDateKey, c)
			return
		}
	} else {
		err = NewCustomEventError("EventDate", c.CustomEventPrototype.EventDateKey, c)
		return
	}
	if c.CustomEventPrototype.AccomplishMinScoreKey == "" {
		event.AccomplishMinScore = 0
	} else {
		intValue, ok = c.Data[c.CustomEventPrototype.AccomplishMinScoreKey].(int)
		if ok {
			event.AccomplishMinScore = intValue
		} else {
			err = NewCustomEventError("AccomplishMinScore", c.CustomEventPrototype.AccomplishMinScoreKey, c)
			return
		}
	}
	if c.CustomEventPrototype.AccomplishMaxScoreKey == "" {
		event.AccomplishMaxScore = 1
	} else {
		intValue, ok = c.Data[c.CustomEventPrototype.AccomplishMaxScoreKey].(int)
		if ok {
			event.AccomplishMaxScore = intValue
		} else {
			err = NewCustomEventError("AccomplishMaxScore", c.CustomEventPrototype.AccomplishMaxScoreKey, c)
			return
		}
	}
	if c.CustomEventPrototype.ExpectedScoreKey == "" {
		event.ExpectedScore = 1
	} else {
		intValue, ok = c.Data[c.CustomEventPrototype.ExpectedScoreKey].(int)
		if ok {
			event.ExpectedScore = intValue
		} else {
			err = NewCustomEventError("ExpectedScore", c.CustomEventPrototype.ExpectedScoreKey, c)
			return
		}
	}
	event.Hook = make(map[string]interface{}, len(c.CustomEventPrototype.HookKeys))
	for i := 0; i < len(c.CustomEventPrototype.HookKeys); i++ {
		val := c.CustomEventPrototype.HookKeys[i]
		event.Hook[val] = c.Data[val]
	}

	stringValue, ok = c.Data[c.CustomEventPrototype.RemindExpirationDateKey].(string)
	if ok {
		dateValue, err = ParseDate(stringValue)
		if err == nil {
			event.RemindInfo.ExpirationDate = &dateValue
		} else {
			err = NewCustomEventError("RemindExpirationDate", c.CustomEventPrototype.RemindExpirationDateKey, c)
			return
		}
	} else {
		err = NewCustomEventError("RemindExpirationDate", c.CustomEventPrototype.RemindExpirationDateKey, c)
		return
	}

	event.RemindInfo.RemindType = c.CustomEventPrototype.RemindTypeKey
	// stringValue, ok = c.Data[c.CustomEventPrototype.RemindTypeKey].(string)
	// if ok {
	// 	event.RemindInfo.RemindType = stringValue
	// } else {
	// 	err = NewCustomEventError("RemindType", c.CustomEventPrototype.RemindTypeKey, c)
	// 	return
	// }

	if c.CustomEventPrototype.RemindMaxScoreKey == "" {
		event.RemindInfo.RemindMaxScore = 1
	} else {
		intValue, ok = c.Data[c.CustomEventPrototype.RemindMaxScoreKey].(int)
		if ok {
			event.RemindInfo.RemindMaxScore = intValue
		} else {
			err = NewCustomEventError("RemindMaxScore", c.CustomEventPrototype.RemindMaxScoreKey, c)
			return
		}
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
		event.RemindHook[val] = c.Data[val]
	}
	return
}

func ParseDate(dateString string) (date time.Time, err error) {
	date, err = time.Parse("2006-01-02T15:04:05Z", dateString)
	if err != nil {
		if !strings.Contains(dateString, "Z") {
			date, err = time.Parse(constants.DateFormatStringYYYYMMDDTHHMMSS, dateString)
			if err != nil {
				if !strings.Contains(dateString, "T") {
					date, err = time.Parse(constants.DateFormatStringYYYYMMDD, dateString)
					if err != nil {
						return
					}
				}
				return
			}
		} else {
			return
		}
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
	err = tmpl.Execute(&buffer, c.Data)
	if err != nil {
		return
	}
	value = buffer.String()
	return
}

// AfterCreate di CustomEvent
func (c *CustomEvent) AfterCreate(tx *gorm.DB) (err error) {
	var event Event
	event, err = c.GetEvent(tx)
	if err != nil {
		return
	}
	return AddEvent(tx, &event)
}

// BeforeDelete di CustomEvent
func (c *CustomEvent) BeforeDelete(tx *gorm.DB) (err error) {
	var event Event
	event, err = c.GetEvent(tx)
	if err != nil {
		return
	}
	return DeleteEvent(tx, &event)
}

// AfterUpdate di CustomEvent
func (c *CustomEvent) AfterUpdate(tx *gorm.DB) (err error) {
	var event Event
	event, err = c.GetEvent(tx)
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
