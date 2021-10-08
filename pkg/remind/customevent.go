package remind

import (
	"bytes"
	"encoding/json"
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
	// un evento può avere associato solo un oggetto/azienda/utente
	ObjectReferenceID string
	// id dell'evento della tabella events che è associato al customEvent
	EventID int
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
	var ok bool
	if c.CustomEventPrototype.ID == 0 {
		// se non c'è già, carico CustomEventPrototype dal db
		if c.CustomEventPrototypeID == 0 {
			// significa che non abbiamo le info necessarie (vedi caso delete in
			// cui valorizziamo per esempio solo ID). Si prendono dal db
			if err = db.Table("custom_events").Preload("CustomEventPrototype").Where("id=?", c.ID).First(&c).Error; err != nil {
				return
			}
		} else if err = db.Table("custom_event_prototypes").Where("id=?", c.CustomEventPrototypeID).First(&c.CustomEventPrototype).Error; err != nil {
			return
		}
	}
	// EventType è id del prototipo in stringa
	event.ID = c.EventID
	event.EventType = strconv.Itoa(c.CustomEventPrototypeID)
	stringValue, ok = c.Data[c.CustomEventPrototype.EventDateKey].(string)
	if ok {
		var dateValue time.Time
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
	event.Hook = make(map[string]interface{}, len(c.CustomEventPrototype.HookKeys)+2)
	event.Hook["object_reference_id"] = c.ObjectReferenceID
	event.Hook["event_type"] = event.EventType
	for i := 0; i < len(c.CustomEventPrototype.HookKeys); i++ {
		val := c.CustomEventPrototype.HookKeys[i]
		event.Hook[val] = c.Data[val]
	}

	var intervalValue models.Interval
	var intervalBytes []byte
	intervalBytes, err = json.Marshal(c.Data[c.CustomEventPrototype.RemindExpirationIntervalKey])
	if err != nil {
		return
	}
	err = json.Unmarshal(intervalBytes, &intervalValue)
	if err != nil {
		return
	}
	dateValue := event.EventDate.AddDate(int(intervalValue.Anni), int(intervalValue.Mesi), intervalValue.Giorni)
	event.RemindInfo.ExpirationDate = &dateValue

	event.RemindInfo.RemindType = c.CustomEventPrototype.RemindTypeKey
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

	event.RemindInfo.ObjectDescription, err = c.getCustomEventReferenceObjectDescription(db)
	if err != nil {
		return
	}
	event.RemindHook = make(map[string]interface{}, len(c.CustomEventPrototype.RemindHookKeys)+2)
	event.RemindHook["object_reference_id"] = c.ObjectReferenceID
	event.RemindHook["event_type"] = event.RemindType
	for i := 0; i < len(c.CustomEventPrototype.RemindHookKeys); i++ {
		val := c.CustomEventPrototype.RemindHookKeys[i]
		event.RemindHook[val] = c.Data[val]
	}
	return
}

func ParseDate(dateString string) (date time.Time, err error) {
	date, err = time.Parse(constants.DateFormatStringYYYYMMDDTHHMMSSZ, dateString)
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

func parseGenericTemplate(values interface{}, templateString string) (value string, err error) {
	var tmpl *template.Template
	tmpl, err = template.New("mind_reminder_template").Parse(templateString)
	if err != nil {
		return
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, values)
	if err != nil {
		return
	}
	value = buffer.String()
	return
}

// parseTemplate parsa templateString prendendo i dati da EventData
// ES: {{.Count}} items are made of {{.Material}}
func (c *CustomEvent) parseTemplate(templateString string) (value string, err error) {
	return parseGenericTemplate(c.Data, templateString)
}

func (c *CustomEvent) getCustomEventReferenceObjectDescription(db *gorm.DB) (description string, err error) {
	// recupero le info della sezione collegata
	var section CustomSection
	if err = db.Where("id=?", c.CustomSectionID).Preload("CustomObjectPrototype").First(&section).Error; err != nil {
		return
	}
	if section.CustomObjectPrototype != nil {
		// si recupera il template della descrizione dal prototipo oggetto
		var obj CustomObject
		var objId int
		objId, err = strconv.Atoi(c.ObjectReferenceID)
		if err != nil {
			return
		}
		if err = db.Where("id=?", objId).First(&obj).Error; err != nil {
			return
		}
		description, err = parseGenericTemplate(obj.Data, section.CustomObjectPrototype.DescriptionTemplate)
		if err != nil {
			return
		}
	} else {
		// si recupera la descrizione da referenceObjectMap
		var referenceObj ReferenceObject
		referenceObj, err = GetReferenceObject(db, section.Reference, c.ObjectReferenceID)
		if err != nil {
			return
		}
		description = (referenceObj).GetDescription()
	}
	return
}

// BeforeCreate di CustomEvent
func (c *CustomEvent) BeforeCreate(tx *gorm.DB) (err error) {
	var event Event
	event, err = c.GetEvent(tx)
	if err != nil {
		return
	}
	if err = AddEvent(tx, &event); err != nil {
		return
	}
	c.EventID = event.ID
	return
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

// BeforeUpdate di CustomEvent
func (c *CustomEvent) BeforeUpdate(tx *gorm.DB) (err error) {
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
