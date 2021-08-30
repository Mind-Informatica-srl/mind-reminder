package mrmodel

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/models"
)

type accomplishers []Accomplisher

func (a accomplishers) Len() int           { return len(a) }
func (a accomplishers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a accomplishers) Less(i, j int) bool { return a[i].AccomplishAt.Before(a[j].AccomplishAt) }

// Remind è la struct delle scadenze
type Remind struct {
	ID int
	// Descrizione della scadenza
	Description *string
	// Tipo della scadenza
	RemindType string
	// json dell'oggetto
	ObjectRaw models.JSONB
	// id dell'oggetto
	ObjectID string
	// tipo della model dell'oggetto
	ObjectType string
	// Data Scadenza
	ExpireAt time.Time
	// Data Creazione
	CreatedAt time.Time `gorm:"default:now()"`
	// Descrizione dello stato della scadenza
	StatusDescription *string
	// criteri di visibilià
	Visibility *string
	// assolvenze
	Accomplishers []Accomplisher `gorm:"foreignKey:remind_id;references:id"`
}

// Accomplished restituisce lo stato di assolvenza della scadenza, la percentuale di assolvimento,
// l'assolvenza determinante e quelle in eccedenza
func (r Remind) Accomplished() (
	accomplished bool,
	percentage float64,
	accomplisher Accomplisher,
	surplus []Accomplisher,
) {
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

// TableName return the remind table name
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

func (r *Remind) MarshalJSON() ([]byte, error) {
	accomplished, percentage, accomplisher, surplus := r.Accomplished()
	return json.Marshal(struct {
		Remind
		Accomplished bool
		Percentage   float64
		Accomplisher Accomplisher
		Surplus      []Accomplisher
	}{
		Remind:       Remind(*r),
		Accomplished: accomplished,
		Percentage:   percentage,
		Accomplisher: accomplisher,
		Surplus:      surplus,
	})
}
