package mindre

import (
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/models"
)

//struct delle scadenze
type Reminder struct {
	ID int
	//Descrizione della scadenza
	Description *string
	//Tipo della scadenza
	ReminderType string
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
	//Percentuale di assolvenza
	Percentage float64
	//Descrizione dello stato della scadenza
	StatusDescription *string
	//criteri di visibilià
	Visibility *string
}

func (t *Reminder) TableName() string {
	return "reminder"
}

// SetPK set the pk for the model
func (t *Reminder) SetPK(pk interface{}) error {
	id := pk.(int)
	t.ID = id
	return nil
}

// VerifyPK check the pk value
func (t *Reminder) VerifyPK(pk interface{}) (bool, error) {
	id := pk.(int)
	return t.ID == id, nil
}
