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
	if db.Statement.Model != nil {
		if isRemindable(db.Statement.Model) {
			_ = addRecord(db, actionCreate)
		}
	}
}

//callback da eseguire dopo after_update
func (p *Plugin) addUpdated(db *gorm.DB) {
	if db.Statement.Model != nil {
		if isRemindable(db.Statement.Model) {
			_ = addUpdateRecord(db, p.opts)
		}
	}
}

//callback da eseguire dopo after_delete.
func (p *Plugin) addDeleted(db *gorm.DB) {
	if isRemindable(db.Statement.Model) {
		_ = addRecord(db, actionDelete)
	}
}

// Writes new reminder row to db.
func addRecord(db *gorm.DB, action string) error {
	cl, err := newToRemind(db, action)
	if err != nil {
		return nil
	}

	return db.Create(cl).Error
}

func addUpdateRecord(db *gorm.DB, opts options) error {
	cl, err := newToRemind(db, actionUpdate)
	if err != nil {
		return err
	}

	return db.Create(cl).Error
}

//restituisce uno slice di scadenze
func newToRemind(db *gorm.DB, action string) (*ToRemind, error) {
	rawObject, err := json.Marshal(db.Statement.Model)
	if err != nil {
		return nil, err
	}
	return &ToRemind{
		Action:     action,
		ObjectID:   interfaceToString(GetPrimaryKeyValue(db)),
		ObjectType: db.Statement.Schema.ModelType.Name(),
		ObjectRaw:  string(rawObject),
	}, nil
}
