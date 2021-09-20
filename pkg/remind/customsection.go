package remind

import "time"

type CustomSection struct {
	ID int
	// titolo della sezione
	Name string
	// descrizione
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	// icona dell asezione
	Icon string
	// ordine di visualizzazione (prima oggetti o prima eventi)
	Configuration SectionConfiguration
	// posizione nel menu
	Order int
	// campo per indicare il riferimento della sezione
	// se risorsa umana, oggetto custom o altro
	Reference string
	// id del prototipo dell'oggetto
	ObjectPrototypeID *int
	// struct CustomObject associato
	ObjectPrototype *CustomObjectPrototype
	// elenco dei CustomEvent
	EventsPrototypes []CustomEvent
}

// SetPK set the pk for the model
func (c *CustomSection) SetPK(pk interface{}) error {
	c.ID = pk.(int)
	return nil
}

// VerifyPK check the pk value
func (c *CustomSection) VerifyPK(pk interface{}) (bool, error) {
	id := pk.(int)
	return c.ID == id, nil
}
