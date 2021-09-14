package remind

import (
	"time"
)

// Accomplisher è un'assolvenza (anche parziale) ad una scadenza
type Accomplisher struct {
	ID           int // anche qui mi fa abbastanza schifo, ma ora mi fa fatica pensarci!
	RemindID     int
	Remind       Remind
	EventID      int
	Event        Event
	AccomplishAt time.Time
	// Percentuale di assolvenza
	Score int
}

// TableName return the accomplishers table name
func (a *Accomplisher) TableName() string {
	return "accomplishers"
}

// SetPK set the pk for the model
func (a *Accomplisher) SetPK(pk interface{}) error {
	id := pk.(int)
	a.ID = id
	return nil
}

// VerifyPK check the pk value
func (a *Accomplisher) VerifyPK(pk interface{}) (bool, error) {
	id := pk.(int)
	return a.ID == id, nil
}

// IsZero return true if the percentage is zero
func (a Accomplisher) IsZero() bool {
	return a.Score == 0
}

type Accomplishers []*Accomplisher

func (a Accomplishers) Len() int           { return len(a) }
func (a Accomplishers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Accomplishers) Less(i, j int) bool { return a[i].AccomplishAt.Before(a[j].AccomplishAt) }
func (accomplishers Accomplishers) Score() (score int) {
	for _, a := range accomplishers {
		score += int(a.Score)
	}
	return
}
