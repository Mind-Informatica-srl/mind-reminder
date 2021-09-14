package remind

import (
	"bytes"
	"text/template"
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/models"
	"gorm.io/gorm"
)

type CustomEvent struct {
	ID                        int
	PrototypeID               int
	EventData                 models.JSONB
	EventTypeKey              string
	EventDateKey              string
	AccomplishMinScoreKey     string
	AccomplishMaxScoreKey     string
	HookKeys                  []string
	RemindExpirationDateKey   string
	RemindTypeKey             string
	RemindMaxScoreKey         string
	RemindDescriptionTemplate string
}

// GetEvent dato EventData in c (*CustomEvent)
// restituisce un evento ed un eventuale errore
func (c *CustomEvent) GetEvent() (event Event, err error) {
	var stringValue string
	var intValue int
	var dateValue time.Time
	var ok bool
	stringValue, ok = c.EventData[c.EventTypeKey].(string)
	if ok {
		event.EventType = stringValue
	}
	dateValue, ok = c.EventData[c.EventDateKey].(time.Time)
	if ok {
		event.EventDate = dateValue
	}
	intValue, ok = c.EventData[c.AccomplishMinScoreKey].(int)
	if ok {
		event.AccomplishMinScore = intValue
	}
	intValue, ok = c.EventData[c.AccomplishMaxScoreKey].(int)
	if ok {
		event.AccomplishMaxScore = intValue
	}
	event.Hook = make(map[string]interface{}, len(c.HookKeys))
	for i := 0; i < len(c.HookKeys); i++ {
		val := c.HookKeys[i]
		event.Hook[val] = c.EventData[val]
	}
	dateValue, ok = c.EventData[c.RemindExpirationDateKey].(time.Time)
	if ok {
		event.Remind.ExpirationDate = dateValue
	}
	stringValue, ok = c.EventData[c.RemindTypeKey].(string)
	if ok {
		event.Remind.RemindType = stringValue
	}
	intValue, ok = c.EventData[c.RemindMaxScoreKey].(int)
	if ok {
		event.Remind.MaxScore = intValue
	}
	stringValue, err = c.parseTemplate(c.RemindDescriptionTemplate)
	if err != nil {
		return
	}
	event.Remind.Description = stringValue
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
	AddEvent(&event)
	return
}

// BeforeDelete di CustomEvent
func (c *CustomEvent) BeforeDelete(tx *gorm.DB) (err error) {
	var event Event
	event, err = c.GetEvent()
	if err != nil {
		return
	}
	DeleteEvent(&event)
	return
}

// AfterUpdate di CustomEvent
func (c *CustomEvent) AfterUpdate(tx *gorm.DB) (err error) {
	var event Event
	event, err = c.GetEvent()
	if err != nil {
		return
	}
	UpdateEvent(&event)
	return
}
