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

// Accomplishers rappresenta un array di accomplisher
type Accomplishers []*Accomplisher

// Len restituisce il numero di elementi
func (accomplishers Accomplishers) Len() int { return len(accomplishers) }

// Swap scambia di posto due elementi
func (accomplishers Accomplishers) Swap(i, j int) {
	accomplishers[i], accomplishers[j] = accomplishers[j], accomplishers[i]
}

// Less restituisce true se l'iesimo accomplisher è precedente al jesimo
func (accomplishers Accomplishers) Less(i, j int) bool {
	return accomplishers[i].AccomplishAt.Before(accomplishers[j].AccomplishAt)
}

// Score restituisce lo score della lista di accomplishers
func (accomplishers Accomplishers) Score() (score int) {
	for _, a := range accomplishers {
		score += a.Score
	}
	return
}

// // Scan implements the database/sql Scanner interface.
// func (dst *Accomplishers) Scan(src interface{}) error {
// 	tmp := pq.GenericArray{
// 		A: dst,
// 	}
// 	fmt.Println(tmp)
// 	if err := tmp.Scan(src); err != nil {
// 		return err
// 	}
// 	*dst = tmp.A.(Accomplishers)
// 	return nil
// }

// // Value implements the database/sql/driver Valuer interface.
// func (src Accomplishers) Value() (driver.Value, error) {
// 	tmp, err := pq.GenericArray{
// 		A: src,
// 	}.Value()
// 	fmt.Println(tmp)
// 	return tmp, err
// }
