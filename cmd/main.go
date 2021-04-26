package main

import (
	"log"
	"strings"

	v1 "github.com/Mind-Informatica-srl/mind-reminder/internal/api/v1"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/models"
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

func (l Utente) Reminders(db *gorm.DB) (toInsert []models.Reminder, toDelete []models.Reminder, err error) {
	toInsert = []models.Reminder{}
	toDelete = []models.Reminder{}
	err = nil
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
