package main

import (
	mindreminder "github.com/Mind-Informatica-srl/mind-reminder"
	"github.com/Mind-Informatica-srl/mind-reminder/config"
	"gorm.io/gorm"
)

func main() {
	db := config.Env.Db
	// mindreminder.Register(db)
	if err := updateRuolo(db); err != nil {
		panic(err)
	}
}

type Ruolo struct {
	ID   int
	Nome string

	*mindreminder.Remind
}

func (r *Ruolo) TableName() string {
	return "ruoli"
}

func (l Ruolo) Reminders(db *gorm.DB) (toInsert []mindreminder.Reminder, toDelete []mindreminder.Reminder, err error) {
	toInsert = []mindreminder.Reminder{}
	toDelete = []mindreminder.Reminder{}
	err = nil
	return
}

func updateRuolo(db *gorm.DB) error {

	var element = Ruolo{
		ID:   1,
		Nome: "Admin",
	}
	if err := db.Save(&element).Error; err != nil {
		return err
	}
	return nil
}
