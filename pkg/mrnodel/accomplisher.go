package mrmodel

import "time"

// Accomplisher è un'assolvenza (anche parziale) ad una scadenza
type Accomplisher struct {
	ID           int // anche qui mi fa abbastanza schifo, ma ora mi fa fatica pensarci!
	RemindID     int
	ObjectID     string
	AccomplishAt time.Time
	//Percentuale di assolvenza
	Percentage float64
}

func (t *Accomplisher) TableName() string {
	return "accomplishers"
}

// SetPK set the pk for the model
func (t *Accomplisher) SetPK(pk interface{}) error {
	id := pk.(int)
	t.ID = id
	return nil
}

// VerifyPK check the pk value
func (t *Accomplisher) VerifyPK(pk interface{}) (bool, error) {
	id := pk.(int)
	return t.ID == id, nil
}

func (a Accomplisher) IsZero() bool {
	return a.Percentage == 0
}
