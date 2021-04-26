package main

import (
	"log"
	"strings"
	"time"

	v1 "github.com/Mind-Informatica-srl/mind-reminder/internal/api/v1"
	mindreminder "github.com/Mind-Informatica-srl/mind-reminder/internal/models"
	"gorm.io/gorm"
)

type Utente struct {
	ID        int
	Username  string
	Nome      string
	Cognome   string
	MasterKey []byte  `gorm:"-" json:"-"`
	JwtToken  *string `gorm:"-"`
	RuoloID   int

	*mindreminder.Remind
}

func (u Utente) LogDescription() string {
	return u.Username
}

func (r *Utente) TableName() string {
	return strings.ToLower("Utenti")
}

func (l *Utente) Reminders(db *gorm.DB) (toInsert []mindreminder.Reminder, toDelete []mindreminder.Reminder, err error) {
	toInsert = []mindreminder.Reminder{}
	toDelete = []mindreminder.Reminder{}
	err = nil
	newEl, err := mindreminder.NewBaseReminder(l, "Test", "Scadenza")
	if err != nil {
		return
	}
	newEl.ExpireAt = time.Now()
	toInsert = append(toInsert, newEl)

	// var oldEl models.Reminder
	// if err = db.Model(oldEl).First(&oldEl).Error; err != nil {
	// 	return
	// }
	oldEl := mindreminder.Reminder{
		ReminderType: "Scadenza",
		ObjectType:   "Utente",
	}
	toDelete = append(toDelete, oldEl)

	return
}

func main() {
	structList := []interface{}{
		Utente{},
	}
	if err := v1.StartService(structList); err != nil {
		log.Fatal(err)
	}

}
