package mrmodel

import (
	"sort"
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/models"
)

type accomplishers []Accomplisher

func (a accomplishers) Len() int           { return len(a) }
func (a accomplishers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a accomplishers) Less(i, j int) bool { return a[i].AccomplishAt.Before(a[j].AccomplishAt) }

//struct delle scadenze
type Remind struct {
	ID int
	//Descrizione della scadenza
	Description *string
	//Tipo della scadenza
	RemindType string
	//json dell'oggetto
	ObjectRaw models.JSONB
	// id dell'oggetto
	ObjectID string
	//tipo della model dell'oggetto
	ObjectType string
	//Data Scadenza
	ExpireAt time.Time
	//Data Assolvenza
	AccomplishedAt *time.Time
	//Data Creazione
	CreatedAt time.Time `gorm:"default:now()"`
	//Descrizione dello stato della scadenza
	StatusDescription *string
	//criteri di visibilià
	Visibility *string
	//assolvenze
	Accomplishers []Accomplisher // TODO: da mettere l'annotazione di gorm
}

// Accomplished restituisce lo stato di assolvenza della scadenza, la percentuale di assolvimento,
// l'assolvenza determinante e quelle in eccedenza
// TODO: da far restituire i campi accomplished, percentage, e volendo accomplisher o accomplisher.AccomplishAt
func (r Remind) Accomplished() (accomplished bool, percentage float64, accomplisher Accomplisher, surplus []Accomplisher) {
	sort.Sort(accomplishers(r.Accomplishers))
	for _, a := range r.Accomplishers {
		percentage += a.Percentage
		if percentage >= 1 && accomplisher.IsZero() {
			accomplisher = a
			accomplished = true
		} else if percentage > 1 {
			surplus = append(surplus, a)
		}
	}
	return
}

// TODO: correggere l'sql di creazione della tabella
func (r *Remind) TableName() string {
	return "remind"
}

// SetPK set the pk for the model
func (r *Remind) SetPK(pk interface{}) error {
	id := pk.(int)
	r.ID = id
	return nil
}

// VerifyPK check the pk value
func (r *Remind) VerifyPK(pk interface{}) (bool, error) {
	id := pk.(int)
	return r.ID == id, nil
}
