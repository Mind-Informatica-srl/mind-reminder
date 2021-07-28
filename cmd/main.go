package main

import (
	"errors"
	"log"
	"strings"
	"time"

	mindreminder "github.com/Mind-Informatica-srl/mind-reminder/pkg"
	models "github.com/Mind-Informatica-srl/mind-reminder/pkg/models"
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

	*models.Remind
}

func (u Utente) LogDescription() string {
	return u.Username
}

func (r *Utente) TableName() string {
	return strings.ToLower("Utenti")
}

func (l *Utente) Reminders(db *gorm.DB, action string) (toInsert []models.Reminder, toDelete []models.Reminder, err error) {
	newEl, err := models.NewBaseReminder(l, "Test", "Scadenza")
	if err != nil {
		return
	}
	newEl.ExpireAt = time.Now()
	toInsert = append(toInsert, newEl)

	var oldEl models.Reminder
	if err = db.Model(oldEl).First(&oldEl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	toDelete = append(toDelete, oldEl)
	return
}

func main() {
	structList := []interface{}{
		Utente{},
	}
	if err := mindreminder.StartService(structList, "test"); err != nil {
		log.Fatal(err)
	}

}
