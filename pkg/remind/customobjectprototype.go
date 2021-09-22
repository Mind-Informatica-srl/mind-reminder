package remind

import "time"

// CustomObjectPrototype prototipo di oggetto di tipo custom
type CustomObjectPrototype struct {
	ID int
	// nome del prototipo
	Name string
	// descrizione del prototipo
	Description *string
	// data di creazione del prototipo
	CreatedAt time.Time
	// data di modifica del prototipo
	UpdatedAt time.Time
	// jsonb con i dati dell'evento
	PrototypeObjectData PrototypeObjectData
	DescriptionTemplate string
}

// SetPK set the pk for the model
func (c *CustomObjectPrototype) SetPK(pk interface{}) error {
	c.ID = pk.(int)
	return nil
}

// VerifyPK check the pk value
func (c *CustomObjectPrototype) VerifyPK(pk interface{}) (bool, error) {
	id := pk.(int)
	return c.ID == id, nil
}
