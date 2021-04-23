package v1

import (
	"reflect"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/models"
	"gorm.io/gorm"
)

type Plugin struct {
	db   *gorm.DB
	opts options
}

//medoto da eseguire nell'init() per registrare le callback
func Register(db *gorm.DB, opts ...Option) (Plugin, error) {
	o := options{}
	o.objectTypes = make(map[string]reflect.Type)
	for _, option := range opts {
		option(&o)
	}
	p := Plugin{db: db, opts: o}

	return p, nil

}

// GetRecords returns all records by objectId.
// Flag prepare allows to decode content of Raw* fields to direct fields, e.g. RawObject to Object.
func (p *Plugin) GetRecords(objectId string, prepare bool) (reminders []models.RemindToCalculate, err error) {
	defer func() {
		if prepare {
			for i := range reminders {
				if t, ok := p.opts.objectTypes[reminders[i].ObjectType]; ok {
					err = reminders[i].PrepareObject(t)
					if err != nil {
						return
					}
				}
			}
		}
	}()
	return reminders, p.db.Where("object_id = ?", objectId).Find(&reminders).Error
}
