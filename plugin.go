package mindreminder

import (
	"reflect"

	"gorm.io/gorm"
)

type Plugin struct {
	db   *gorm.DB
	opts options
}

//medoto da eseguire nell'init() per registrare le callback
func Register(db *gorm.DB, opts ...Option) (Plugin, error) {
	//crea su db automaticamente la tabella per reminder (se non esiste)
	// err := db.AutoMigrate(&ToRemind{})
	// if err != nil {
	// 	return Plugin{}, err
	// }
	o := options{}
	o.objectTypes = make(map[string]reflect.Type)
	for _, option := range opts {
		option(&o)
	}
	p := Plugin{db: db, opts: o}
	callback := db.Callback()
	//si registrano le callback
	callback.Create().After("gorm:after_create").Register("mindreminder:create", p.addCreated)
	callback.Update().After("gorm:after_update").Register("mindreminder:update", p.addUpdated)
	callback.Delete().After("gorm:after_delete").Register("mindreminder:delete", p.addDeleted)

	return p, nil

}

// GetRecords returns all records by objectId.
// Flag prepare allows to decode content of Raw* fields to direct fields, e.g. RawObject to Object.
func (p *Plugin) GetRecords(objectId string, prepare bool) (reminders []ToRemind, err error) {
	defer func() {
		if prepare {
			for i := range reminders {
				if t, ok := p.opts.objectTypes[reminders[i].ObjectType]; ok {
					err = reminders[i].prepareObject(t)
					if err != nil {
						return
					}
				}
			}
		}
	}()
	return reminders, p.db.Where("object_id = ?", objectId).Find(&reminders).Error
}
