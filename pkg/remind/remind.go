package remind

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	"github.com/Mind-Informatica-srl/restapi/pkg/actions"
	"github.com/Mind-Informatica-srl/restapi/pkg/controllers"
	"github.com/Mind-Informatica-srl/restapi/pkg/models"
	"gorm.io/gorm"
)

// Remind è la struct delle scadenze
type Remind struct {
	ID int
	// Descrizione della scadenza
	Description *string
	// Tipo della scadenza
	RemindType string
	// id dell'evento
	EventID int
	Event   Event
	// Data Scadenza
	ExpireAt time.Time
	// Data Creazione
	CreatedAt time.Time `gorm:"default:now()"`
	// Descrizione dello stato della scadenza
	StatusDescription *string
	// criteri di visibilià
	Visibility *string
	// assolvenze
	accomplishers Accomplishers `gorm:"foreignKey:remind_id;references:id"`
	MaxScore      int
}

// accomplished restituisce lo stato di assolvenza della scadenza, la percentuale di assolvimento,
// l'assolvenza determinante e quelle in eccedenza
func (r Remind) accomplished() (
	accomplished bool,
	score int,
	accomplisher *Accomplisher,
	surplus []Accomplisher,
) {
	sort.Sort(r.accomplishers)
	for _, a := range r.accomplishers {
		score += a.Score
		if score >= 1 && accomplisher.IsZero() {
			accomplisher = a
			accomplished = true
		} else if score > 1 {
			surplus = append(surplus, *a)
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

// MarshalJSON produce il json del remind con le informazioni utili
func (r *Remind) MarshalJSON() ([]byte, error) {
	accomplished, score, accomplisher, surplus := r.accomplished()
	return json.Marshal(struct {
		Remind
		Accomplished bool
		Score        int
		Accomplisher *Accomplisher
		Surplus      []Accomplisher
	}{
		Remind:       *r,
		Accomplished: accomplished,
		Score:        score,
		Accomplisher: accomplisher,
		Surplus:      surplus,
	})
}

// AfterCreate cerca di assolvere il remind appena inserito
func (r *Remind) AfterCreate(tx *gorm.DB) (err error) {
	// recupera gli eventi riguardanti stessa azienda, unità locale, utente e tipo
	// successivi a remind, che assolvono remind successivi
	err = r.searchForAccomplishers(tx)
	return
}

func (r *Remind) searchForAccomplishers(tx *gorm.DB) (err error) {
	for {
		var event Event
		if err = tx.Joins("(select sum(score) as tot_score, max(accomplish_at) as max_date, event_id "+
			"from accomplishers group by event_id) as accstatus on accstatus.event_id = event.id").
			Where("accstatus.tot_score < event.accomplish_max_score or max_date > ?", r.Event.EventDate).
			Where("event_type = ? and hook = ? and event_date > ? and event_date <= ?", r.RemindType, r.Event.Hook, r.Event.EventDate, r.ExpireAt).
			Order("event.event_date").
			Preload("Accomplishers.Event").
			First(&event).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = nil
			}
			return
		}
		// creo l'accomplisher
		a := createAccomplisher(event, *r)
		if err = tx.Create(&a).Error; err != nil {
			return
		}
		if err = event.tryToAccomplish(tx); err != nil {
			return
		}
		r.accomplishers = append(r.accomplishers, &a)
		if r.accomplishers.Score() >= r.MaxScore {
			return
		}
	}
}

var remindDelegate = models.NewBaseDelegate(func() *gorm.DB {
	return config.Current().DB
}, func() models.PKModel {
	return &Remind{}
},
	func() interface{} {
		return &[]Remind{}
	}, func(r *http.Request) (interface{}, error) {
		return actions.PrimaryKeyIntExtractor(r, "ID")
	})

// RemindController è il controller dei remind
var RemindController = controllers.CreateModelController("/remind", remindDelegate)

// Create create a configuration
func Create(db *gorm.DB) (*config.Data, error) {
	config.CurrentConfiguration = &config.Data{
		DB: db,
	}

	return config.CurrentConfiguration, nil
}
