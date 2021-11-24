package remind

import (
	"errors"
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/models"
	"gorm.io/gorm"
)

// type RemindStrategy string

// const (
// 	Always         RemindStrategy = ""
// 	MaxScore       RemindStrategy = "max_score"
// 	ZeroOrMaxScore RemindStrategy = "zero_or_max_score"
// )

type RemindInfo struct {
	ExpirationDate    *time.Time
	RemindType        string
	RemindMaxScore    int
	RemindDescription string
	ObjectDescription string
	RemindHook        models.JSONB
}

// Event rappresenta un evento che può generare un remind e assolverne altri
type Event struct {
	ID                 int
	EventType          string
	EventDate          *time.Time
	AccomplishMinScore int
	AccomplishMaxScore int
	ExpectedScore      int
	Accomplishers      Accomplishers `gorm:"foreignKey:event_id;references:id"`
	Hook               models.JSONB
	RemindInfo         `gorm:"embedded"`
}

// AfterCreate cerca le scadenze a cui assolve l'evento inserito ed eventualmente genera la scadenza
func (e *Event) AfterCreate(tx *gorm.DB) (err error) {
	return e.elaborateEvent(tx)
}

// BeforeDelete elimina le assolvenze e la scadenza generati dall'evento cancellato
func (e *Event) BeforeDelete(tx *gorm.DB) (err error) {
	// elimimo le assolvenze e rivaluto i remind relativi
	var accs Accomplishers
	if err = tx.Where("event_id = ?", e.ID).Find(&accs).Error; err != nil {
		return
	}
	if accs.Len() > 0 {
		// where 1=1 per non avere gorm.ErrMissingWhereClause
		if err = tx.Where("1 = 1").Delete(&accs).Error; err != nil {
			return
		}

	}
	// elimino l'eventuale remind (e faccio in modo che tutti gli eventi che lo assolvevamo vengano rivalutati)
	if err = tx.Where("event_id = ?", e.ID).Delete(&Remind{}).Error; err != nil {
		return
	}

	for i := range accs {
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
	// recupero il remind con le assolvenze
	var remind Remind
	if err = tx.Where("event_id = ?", e.ID).
		Preload("Accomplishers.Event.Accomplishers").
		First(&remind).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if remind.ID > 0 {
		// si eliminano prima le assolvenze
		if remind.Accomplishers.Len() > 0 {
			// where 1=1 per non avere gorm.ErrMissingWhereClause
			if err = tx.Where("1 = 1").Delete(&remind.Accomplishers).Error; err != nil {
				return
			}
		}
		// elimino il remind
		if err = tx.Delete(&remind).Error; err != nil {
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
			if err = a.Event.elaborateEvent(tx); err != nil {
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

func (e *Event) getRemindFromEvent() Remind {
	return Remind{
		RemindDescription: &e.RemindInfo.RemindDescription,
		ObjectDescription: &e.RemindInfo.ObjectDescription,
		RemindType:        e.RemindInfo.RemindType,
		ExpireAt:          *e.RemindInfo.ExpirationDate,
		CreatedAt:         time.Now(),
		MaxScore:          e.RemindInfo.RemindMaxScore,
		EventID:           e.ID,
		Hook:              e.RemindHook,
	}
}

// inserisce un nuovo Remind prendendo i dati da "e" di tipo Event
func (e *Event) generateRemind(tx *gorm.DB) error {
	if e.RemindInfo.ExpirationDate != nil {
		// se c'è ExpirationDate si genera remind
		remind := e.getRemindFromEvent()
		return tx.Create(&remind).Error
	}
	return nil
}

func (e *Event) elaborateEvent(tx *gorm.DB) (err error) {
	// cerco i remind che posso assolvere e creo le assolvenze
	var hasToGenerateRemind bool
	hasToGenerateRemind, err = e.tryToAccomplish(tx)
	if err != nil {
		return
	}
	// se non ho raggiunto il minimo do errore
	if e.Accomplishers.Score() < e.AccomplishMinScore {
		err = errors.New("min score not reached")
		return
	}

	// switch e.RemindStrategy {
	// case MaxScore:
	// 	if e.Accomplishers.Score() < e.AccomplishMaxScore {
	// 		return
	// 	}
	// case ZeroOrMaxScore:
	// 	if e.AccomplishMaxScore == 0 || e.Accomplishers.Score() < e.AccomplishMaxScore {
	// 		return
	// 	}
	// }
	if hasToGenerateRemind {
		// inserisco il remind
		err = e.generateRemind(tx)
		if err != nil {
			return
		}
	}
	return
}

// addRemindsFromNonAccomplishedEvents inserisce eventuali Remind
// generandoli da quegli eventi con stesso hook e tipo di "e"
// e il cui AccomplishMaxScore sommato a quello degli eventi precedenti
// raggiunge ExpectedScore di "e"
func (e *Event) addRemindsFromNonAccomplishedEvents(tx *gorm.DB) (err error) {
	var events []Event
	events, err = e.getMaxScoreEvents(tx)
	if err != nil {
		return
	}
	for _, v := range events {
		err = v.generateRemind(tx)
		if err != nil {
			return
		}

	}
	return
}

func (e *Event) tryToAccomplish(tx *gorm.DB) (hasToGenerateRemind bool, err error) {
	for {
		// finché c'è un remind da assolvere
		var remind Remind
		if err = e.searchForFirstRemind(tx, &remind); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = e.addRemindsFromNonAccomplishedEvents(tx); err != nil {
					return
				}
				err = nil
			}
			return
		}
		a := createAccomplisher(*e, remind)
		if err = tx.Create(&a).Error; err != nil {
			return
		}
		// inserisco "a" dentro lo slice di remind.Accomplishers (secondo l'ordine delle date dell'evento)
		newEventDate := a.Event.EventDate
		var index int
		// si cerca l'indice in cui dover inserire "a"
		if remind.Accomplishers.Len() > 0 {
			for i := 0; i < len(remind.Accomplishers); i++ {
				if newEventDate.After(*remind.Accomplishers[i].Event.EventDate) {
					index = i
					break
				}
			}
		}
		// remind.Accomplishers = append(remind.Accomplishers, &a)
		remind.Accomplishers = remind.Accomplishers.Insert(index, &a)

		// valuto il remind e tratto il surplus
		_, _, finalAccomplisher, surplus := remind.Accomplished()
		if finalAccomplisher != nil && finalAccomplisher.EventID == e.ID {
			for _, v := range remind.Accomplishers {
				if v.ID == finalAccomplisher.ID && v.Score != finalAccomplisher.Score {
					// si aggiorna l'accomplisher (perchè è stato modificato il suo score)
					if err = tx.Save(&finalAccomplisher).Error; err != nil {
						return
					}
				}
			}
			hasToGenerateRemind = true
		}
		for i := range surplus {
			// elimino il surplus
			if surplus[i].ID > 0 {
				if err = tx.Delete(&surplus[i]).Error; err != nil {
					return
				}
			}
			// controllo l'evento
			var event Event
			if err = tx.Where("ID = ?", surplus[i].EventID).Preload("Accomplishers").First(&event).Error; err != nil {
				return
			}
			// aggiorno le assolvenze
			if err = event.elaborateEvent(tx); err != nil {
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
func (e *Event) searchForFirstRemind(tx *gorm.DB, remind *Remind) (err error) {
	return tx.Joins("Event", tx.Where("event_date <= ?", e.EventDate)).
		Joins("left join (select sum(score) as tot_score, max(accomplish_at) as max_date, remind_id "+
			"from accomplishers group by remind_id) as accstatus on accstatus.remind_id = remind.id").
		Where("coalesce(accstatus.tot_score,0) < remind.max_score or max_date > ?", e.EventDate).
		Where("\"Event\".remind_type = ? and remind.hook = ? and expire_at >= ?", e.EventType, e.Hook, e.EventDate).
		Order("\"Event\".event_date").
		Preload("Accomplishers.Event").
		First(remind).Error
}

func createAccomplisher(event Event, remind Remind) (a Accomplisher) {
	// score := remind.Accomplishers.Score()
	// delta := int(math.Min(float64(event.AccomplishMaxScore-event.Accomplishers.Score()), float64(remind.MaxScore-score)))
	a = Accomplisher{
		RemindID:     remind.ID,
		EventID:      event.ID,
		AccomplishAt: *event.EventDate,
		Score:        event.AccomplishMaxScore - event.Accomplishers.Score(),
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

// getMaxScoreEvents restituisce gli eventi dello stesso tipo il cui AccomplishMaxScore
// sommato a quello dei precedenti in ordine di data raggiunge ExpectedScore di event
func (e *Event) getMaxScoreEvents(db *gorm.DB) (events []Event, err error) {
	var acc Accomplisher
	var list []Event
	err = db.Select("events.*").Joins("left join "+acc.TableName()+" as acc on acc.event_id=events.id").
		Where("hook = ? and remind_type = ?", e.Hook, e.RemindType).Order("event_date").
		Where("acc.id is null").
		Where("event_date <= ?", e.EventDate).
		Find(&list).Error
	if err != nil {
		return
	}
	score := 0
	for i := 0; i < len(list); i++ {
		score += list[i].AccomplishMaxScore
		if score >= e.ExpectedScore {
			score = 0
			events = append(events, list[i])
		}
	}
	return
}
