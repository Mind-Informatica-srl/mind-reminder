package remind

import (
	"errors"
	"math"
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/models"
	"gorm.io/gorm"
)

type RemindStrategy string

const (
	Always         RemindStrategy = "always"
	MaxScore       RemindStrategy = "max_score"
	ZeroOrMaxScore RemindStrategy = "zero_or_max_score"
)

type RemindInfo struct {
	ExpirationDate    *time.Time
	RemindType        string
	MaxScore          int
	RemindDescription string
	ObjectDescription string
}

// Event rappresenta un evento che può generare un remind e assolverne altri
type Event struct {
	ID                 int
	EventType          string
	EventDate          *time.Time
	AccomplishMinScore int
	AccomplishMaxScore int
	RemindStrategy     RemindStrategy
	Accomplishers      Accomplishers `gorm:"foreignKey:event_id;references:id"`
	Hook               models.JSONB
	RemindInfo         `gorm:"embedded"`
}

// AfterCreate cerca le scadenze a cui assolve l'evento inserito ed eventualmente genera la scadenza
func (e *Event) AfterCreate(tx *gorm.DB) (err error) {
	// cerco i remind che posso assolvere e creo le assolvenze
	if err = e.tryToAccomplish(tx); err != nil {
		return
	}
	// se non ho raggiunto il minimo do errore
	if e.Accomplishers.Score() < e.AccomplishMinScore {
		err = errors.New("min score not reached")
		return
	}

	switch e.RemindStrategy {
	case MaxScore:
		if e.Accomplishers.Score() < e.AccomplishMaxScore {
			return
		}
	case ZeroOrMaxScore:
		if e.AccomplishMaxScore == 0 || e.Accomplishers.Score() < e.AccomplishMaxScore {
			return
		}
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
		if err = tx.Where("id = ?", accs[i].RemindID).
			Preload("Accomplishers").First(&remind).Error; err != nil {
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
	var remind Remind
	if err = tx.Where("event_id = ?", e.ID).
		Preload("Accomplishers.Event.Accomplishers").
		First(&remind).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if remind.ID > 0 {
		// lo elimino
		if err = tx.Delete(&remind).Error; err != nil {
			return
		}
	}
	// recupero le assolvenze ad altri remind
	var accs Accomplishers
	if err = tx.Where("event_id = ?", e.ID).Find(&accs).Error; err != nil {
		return
	}
	// le elimino
	for i := range accs {
		if err = tx.Delete(&accs[i]).Error; err != nil {
			return
		}
	}
	// chiamo AfterCreate
	if err = e.AfterCreate(tx); err != nil {
		return
	}
	if remind.ID > 0 {
		// per ogni assolvenza del remind, controllo l'evento relativo
		for _, a := range remind.Accomplishers {
			if err = a.Event.tryToAccomplish(tx); err != nil {
				return
			}
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
		RemindDescription: &e.RemindInfo.RemindDescription,
		ObjectDescription: &e.RemindInfo.ObjectDescription,
		RemindType:        e.RemindInfo.RemindType,
		ExpireAt:          *e.RemindInfo.ExpirationDate,
		CreatedAt:         time.Now(),
		MaxScore:          e.RemindInfo.MaxScore,
		EventID:           e.ID,
	}
}

func (e *Event) tryToAccomplish(tx *gorm.DB) (err error) {
	for {
		// finché c'è un remind da assolvere
		var remind Remind
		if err = e.searchForFirstRemind(tx, &remind); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = nil
			}
			return
		}
		a := createAccomplisher(*e, remind)
		if err = tx.Create(&a).Error; err != nil {
			return
		}
		// valuto il remind e tratto il surplus
		remind.Accomplishers = append(remind.Accomplishers, &a)
		_, _, _, surplus := remind.Accomplished()
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
		e.Accomplishers = append(e.Accomplishers, &a)
		if e.Accomplishers.Score() >= e.AccomplishMaxScore {
			return
		}
	}
}

// searchForFirstRemind seleziona il primo remind in ordine di data con tipo e hook uguale all'evento
// con assoluzioni posteriori all'evento o che ancora deve essere completamente assolto
func (e Event) searchForFirstRemind(tx *gorm.DB, remind *Remind) (err error) {
	return tx.Joins("Event", tx.Where("event_date <= ?", e.EventDate)).
		Joins("left join (select sum(score) as tot_score, max(accomplish_at) as max_date, remind_id "+
			"from accomplishers group by remind_id) as accstatus on accstatus.remind_id = remind.id").
		Where("coalesce(accstatus.tot_score,0) < remind.max_score or max_date > ?", e.EventDate).
		Where("\"Event\".remind_type = ? and hook = ? and expire_at >= ?", e.EventType, e.Hook, e.EventDate).
		Order("\"Event\".event_date").
		Preload("Accomplishers.Event").
		First(remind).Error
}

func createAccomplisher(event Event, remind Remind) (a Accomplisher) {
	score := remind.Accomplishers.Score()
	delta := int(math.Min(float64(event.AccomplishMaxScore-event.Accomplishers.Score()), float64(remind.MaxScore-score)))
	a = Accomplisher{
		RemindID:     remind.ID,
		EventID:      event.ID,
		AccomplishAt: *event.EventDate,
		Score:        delta,
	}
	return
}

// AddEvent aggiunge un evento
func AddEvent(db *gorm.DB, event *Event) (err error) {
	err = db.Create(event).Error
	return
}

// UpdateEvent modifica un evento
func UpdateEvent(db *gorm.DB, event *Event) (err error) {
	err = db.Save(event).Error
	return
}

// DeleteEvent elimina un evento
func DeleteEvent(db *gorm.DB, event *Event) (err error) {
	err = db.Delete(event).Error
	return
}
