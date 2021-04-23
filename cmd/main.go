package main

import (
	"log"
	"strings"
	"time"

	mindlogger "github.com/Mind-Informatica-srl/mind-logger"
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
	mindlogger.CreateLogFolder()
	appLog := mindlogger.CreateLogger()
	timeStart := time.Now()
	appLog.AppendLn("SERVICE START: " + timeStart.Format("01-02-2006"))
	structList := []interface{}{
		Utente{},
	}
	v1.RegisterTypes(structList)
	if err := v1.RicalcolaScadenze(appLog); err != nil {
		appLog.Prepend(err.Error())
		appLog.PrependLn("ERROR")
		timeEnd := time.Now()
		appLog.AppendLn("SERVICE END: " + timeEnd.Format("01-02-2006"))
		appLog.WriteLog()
		log.Fatal(err)
	}
	timeEnd := time.Now()
	appLog.AppendLn("SERVICE END: " + timeEnd.Format("01-02-2006"))
	appLog.WriteLog()
}
