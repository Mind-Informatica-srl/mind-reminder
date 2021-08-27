package mrmodel

import "time"

// Accomplisher è un'assolvenza (anche parziale) ad una scadenza
type Accomplisher struct {
	ID           int // anche qui mi fa abbastanza schifo, ma ora mi fa fatica pensarci!
	RemindID     int
	ObjectID     string
	AccomplishAt time.Time
	// Percentuale di assolvenza
	Percentage float64
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
	return a.Percentage == 0
}
