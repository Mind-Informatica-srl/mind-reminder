package mindreminder

import (
	"encoding/json"

	"gorm.io/gorm"
)

const (
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"
)

//callback da eseguire dopo after_create
func (p *Plugin) addCreated(db *gorm.DB) {
	if db.Statement.Model != nil && isRemindable(db.Statement.Model) {
		_ = addRecord(db, actionCreate)
	}
}

//callback da eseguire dopo after_update
func (p *Plugin) addUpdated(db *gorm.DB) {
	if db.Statement.Model != nil && isRemindable(db.Statement.Model) {
		_ = addRecord(db, actionUpdate)
	}
}

//callback da eseguire dopo after_delete.
func (p *Plugin) addDeleted(db *gorm.DB) {
	if db.Statement.Model != nil && isRemindable(db.Statement.Model) {
		_ = addRecord(db, actionDelete)
	}
}

// Writes new reminder row to db.
func addRecord(db *gorm.DB, action string) error {
	r, err := newToRemind(db, action)
	if err != nil {
		return nil
	}
	return db.Model(&r).Create(&r).Error
}

//restituisce uno slice di scadenze
func newToRemind(db *gorm.DB, action string) (ToRemind, error) {
	rawObject, err := json.Marshal(db.Statement.Model)
	if err != nil {
		return ToRemind{}, err
	}
	return ToRemind{
		Action:     action,
		ObjectID:   interfaceToString(GetPrimaryKeyValue(db)),
		ObjectType: db.Statement.Schema.ModelType.Name(),
		ObjectRaw:  string(rawObject),
	}, nil
}
