package remind

import (
	"errors"
	"math"
	"time"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	"github.com/Mind-Informatica-srl/restapi/pkg/models"
	"gorm.io/gorm"
)

// Event rappresenta un evento che può generare un remind e assolverne altri
type Event struct {
	ID                 int
	EventType          string
	EventDate          time.Time
	AccomplishQuery    *string
	AccomplishMinScore int
	AccomplishMaxScore int
	accomplishers      Accomplishers `gorm:"foreignKey:event_id;references:id"`
	Hook               models.JSONB
	Remind             struct {
		ExpirationDate time.Time
		RemindType     string
		MaxScore       int
		Description    string
	}
}

// AfterCreate cerca le scadenze a cui assolve l'evento inserito ed eventualmente genera la scadenza
func (e *Event) AfterCreate(tx *gorm.DB) (err error) {
	// cerco i remind che posso assolvere e creo le assolvenze
	if err = e.tryToAccomplish(tx); err != nil {
		return
	}
	// se non ho raggiunto il minimo do errore
	if e.accomplishers.Score() < e.AccomplishMinScore {
		err = errors.New("min score not reached")
		return
	}
	// se non ho raggiunto il massimo esco senza generare la scadenza
	if e.accomplishers.Score() < e.AccomplishMaxScore {
		return
	}
	// inserisco il remind
	remind := e.generateRemind()
	err = tx.Create(&remind).Error
	return
}

// BeforeDelete elimina le assolvenze e la scadenza generati dall'evento cancellato
func (e *Event) BeforeDelete(tx *gorm.DB) (err error) {
	// elimino l'eventuale remind (e faccio in modo che tutti gli eventi che lo assolvevamo vengano rivalutati)
	if err = tx.Where("event_id = ?", e.ID).Delete(&Remind{}).Error; err != nil {
		return
	}
	// elimimo le assolvenze e rivaluto i remind relativi
	var accs Accomplishers
	if err = tx.Where("event_id = ?", e.ID).Find(&accs).Error; err != nil {
		return
	}
	for i := range accs {
		if err = tx.Delete(&accs[i]).Error; err != nil {
			return
		}
		var remind *Remind
		if err = tx.Where("id = ?", accs[i].RemindID).Preload("Accomplishers").First(&remind).Error; err != nil {
			return
		}
		if err = remind.searchForAccomplishers(tx); err != nil {
			return
		}
	}
	return
}

// AfterUpdate ripristina le assolvenze e la scadenza dell'evento dopo le modifiche
func (e *Event) AfterUpdate(tx *gorm.DB) (err error) {
	// recupero il remind con le assolvenze
	var remind *Remind
	if err = tx.Where("event_id = ?", e.ID).Preload("Accomplishers.Event.Accomplishers").First(&remind).Error; err != nil {
		return
	}
	// lo elimino
	if err = tx.Delete(&remind).Error; err != nil {
		return
	}
	// recupero le assolvenze ad altri remind
	var accs Accomplishers
	if err = tx.Where("event_id = ?", e.ID).Preload("Event.Remind").Find(&accs).Error; err != nil {
		return
	}
	// le elimino
	for i := range accs {
		if err = tx.Delete(&accs[i]).Error; err != nil {
			return
		}
	}
	// chiamo beforeCreate
	if err = e.AfterCreate(tx); err != nil {
		return
	}
	// per ogni assolvenza del remind, controllo l'evento relativo
	for _, a := range remind.accomplishers {
		if err = a.Event.tryToAccomplish(tx); err != nil {
			return
		}
	}
	// per ogni remind prima assolto, lo controllo
	for _, a := range accs {
		if err = a.Remind.searchForAccomplishers(tx); err != nil {
			return
		}
	}
	return nil
}

func (e Event) generateRemind() (remind Remind) {
	return Remind{
		Description: &e.Remind.Description,
		RemindType:  e.Remind.RemindType,
		ExpireAt:    e.Remind.ExpirationDate,
		CreatedAt:   time.Now(),
		MaxScore:    e.Remind.MaxScore,
		EventID:     e.ID,
	}
}

func (e *Event) tryToAccomplish(tx *gorm.DB) (err error) {
	for {
		// finché c'è un remind da assolvere
		var remind *Remind
		if err = e.searchForFirstRemind(tx, remind); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = nil
			}
			return
		}
		a := createAccomplisher(*e, *remind)
		if err = tx.Create(&a).Error; err != nil {
			return
		}
		// valuto il remind e tratto il surplus
		remind.accomplishers = append(remind.accomplishers, &a)
		_, _, _, surplus := remind.accomplished()
		for i := range surplus {
			// elimino il surplus
			if err = tx.Delete(&surplus[i]).Error; err != nil {
				return
			}
			// controllo l'evento
			var event Event
			if err = tx.Where("ID = ?", a.EventID).Preload("Accomplishers").First(&event).Error; err != nil {
				return
			}
			// aggiorno le assolvenze
			if err = event.tryToAccomplish(tx); err != nil {
				return
			}
		}

		// se ho raggiunto il massimo mi fermo
		e.accomplishers = append(e.accomplishers, &a)
		if e.accomplishers.Score() >= e.AccomplishMaxScore {
			return
		}
	}
}

// searchForFirstRemind seleziona il primo remind in ordine di data con tipo e hook uguale all'evento
// con assoluzioni posteriori all'evento o che ancora deve essere completamente assolto
func (e Event) searchForFirstRemind(tx *gorm.DB, remind *Remind) (err error) {
	return tx.Joins("Event", tx.Where("event_date < ?,", e.EventDate)).
		Joins("(select sum(score) as tot_score, max(accomplish_at) as max_date, remind_id "+
			"from accomplishers group by remind_id) as accstatus on accstatus.remind_id = remind.id").
		Where("accstatus.tot_score < remind.max_score or max_date > ?", e.EventDate).
		Where("remind_type = ? and hook = ?", e.EventType, e.Hook).
		Order("event.event_date").
		Preload("Accomplishers.Event").
		First(&remind).Error
}

func createAccomplisher(event Event, remind Remind) (a Accomplisher) {
	score := remind.accomplishers.Score()
	delta := int(math.Min(float64(event.AccomplishMaxScore-event.accomplishers.Score()), float64(remind.MaxScore-score)))
	a = Accomplisher{
		RemindID:     remind.ID,
		EventID:      event.ID,
		AccomplishAt: event.EventDate,
		Score:        delta,
	}
	return
}

// AddEvent aggiunge un evento
func AddEvent(event *Event) (err error) {
	db := config.Current().DB
	err = db.Create(event).Error
	return
}

// UpdateEvent modifica un evento
func UpdateEvent(event *Event) (err error) {
	db := config.Current().DB
	err = db.Save(event).Error
	return
}

// DeleteEvent elimina un evento
func DeleteEvent(event *Event) (err error) {
	db := config.Current().DB
	err = db.Delete(event).Error
	return
}
