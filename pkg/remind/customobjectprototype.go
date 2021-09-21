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
