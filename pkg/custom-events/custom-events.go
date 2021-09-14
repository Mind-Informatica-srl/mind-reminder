package customevents

import (
	"bytes"
	"text/template"
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/models"
	"gorm.io/gorm"
)

type Strategy struct {
	*models.JSONB
	Reminder struct {
		DescriptionTemplate string
		ExpirationDate      struct {
			FixedDate             *time.Time
			CustomDateKey         string
			IntervalFromDateValue *time.Duration
			IntervalFromDateKey   string
		}
		MinScoreKey string
		MaxScoreKey string
	}
	Accomplisher struct {
		Query             string
		MaxScoreKey       string
		MinScoreKey       string
		AccomplishDateKey string
		AccomplisherValue *int
		PointsKey         string
	}
}

type CustomEvent struct {
	*gorm.Model
	PropertyOfID         int
	CreatedByID          int
	AziendaTargetID      *int
	UnitaLocaleTargetID  *int
	PersonaTargetID      *int
	CustomObjectTargetID *int
	ProototypeID         *int
	Type                 string
	EventData            models.JSONB
	EventDateKey         string
	IdentifierTemplate   string
	Strategy             Strategy
	NextEvent            func(tx *gorm.DB) (nextEvent *CustomEvent, err error)
}

func (e CustomEvent) accMaxScore() (maxScore int) {
	maxScore = e.EventData[e.Strategy.Accomplisher.MaxScoreKey].(int)
	return
}

func (e CustomEvent) accMinScore() (maxScore int) {
	maxScore = e.EventData[e.Strategy.Accomplisher.MinScoreKey].(int)
	return
}

func (e CustomEvent) eventIdentifier() (ID string, err error) {
	var t *template.Template
	if t, err = template.New("ID").Parse(e.IdentifierTemplate); err != nil {
		return
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, e); err != nil {
		return
	}
	ID = buf.String()
	return
}
