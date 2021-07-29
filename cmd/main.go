package main

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/calc"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/model"
	"github.com/Mind-Informatica-srl/mind-reminder/pkg/mindre"
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

	*mindre.Remind
}

func (u Utente) LogDescription() string {
	return u.Username
}

func (r *Utente) TableName() string {
	return strings.ToLower("Utenti")
}

func (l *Utente) Reminders(db *gorm.DB, action string) (toInsert []model.Reminder, toDelete []model.Reminder, err error) {
	newEl, err := calc.NewBaseReminder(l, "Test", "Scadenza")
	if err != nil {
		return
	}
	newEl.ExpireAt = time.Now()
	toInsert = append(toInsert, model.Reminder(newEl))

	var oldEl model.Reminder
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
	if err := mindre.StartService(structList, "test"); err != nil {
		log.Fatal(err)
	}

}
