package models

import "time"

//struct delle scadenze
type Reminder struct {
	ID int
	//Descrizione della scadenza
	Description *string
	//Tipo della scadenza
	ReminderType string
	//json dell'oggetto
	ObjectRaw string `gorm:"type:json"`
	//tipo della model dell'oggetto
	ObjectType string
	//Data Scadenza
	ExpireAt time.Time
	//Data Assolvenza
	AccomplishedAt *time.Time
	//Data Creazione
	CreatedAt time.Time `gorm:"default:now()"`
	//Percentuale di assolvenza
	Percentage float64
	//Descrizione dello stato della scadenza
	StatusDescription *string
	//criteri di visibilià
	Visibility *string
}

func (t *Reminder) TableName() string {
	return "reminder"
}
