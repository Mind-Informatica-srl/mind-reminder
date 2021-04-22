package mindreminder

import (
	"gorm.io/gorm"
)

const (
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"
)

type UpdateDiff map[string]interface{}

//callback da eseguire dopo after_create
func (p *Plugin) addCreated(db *gorm.DB) {
	if db.Statement.Schema != nil {
		if isRemindable(db.Statement.Schema) {
			_ = addRecord(db, actionCreate)
		}
	}
}

//callback da eseguire dopo after_update
func (p *Plugin) addUpdated(db *gorm.DB) {
	if db.Statement.Schema != nil {
		if isRemindable(db.Statement.Schema) {
			_ = addUpdateRecord(db, p.opts)
		}
	}
}

//callback da eseguire dopo after_delete.
func (p *Plugin) addDeleted(db *gorm.DB) {
	if isRemindable(db.Statement.Schema) {
		_ = addRecord(db, actionDelete)
	}
}

// Writes new reminder row to db.
func addRecord(db *gorm.DB, action string) error {
	cl, err := newReminder(db, action)
	if err != nil {
		return nil
	}

	return db.Create(cl).Error
}

func addUpdateRecord(db *gorm.DB, opts options) error {
	cl, err := newReminder(db, actionUpdate)
	if err != nil {
		return err
	}

	return db.Create(cl).Error
}

//restituisce uno slice di scadenze
func newReminder(db *gorm.DB, action string) (*ToRemind, error) {
	// rawObject, err := json.Marshal(db.Statement.Schema)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
	// return &ToRemind{
	// 	ObjectID:   interfaceToString(db.Statement.Schema.PrimaryFields()),
	// 	ObjectType: db.GetModelStruct().ModelType.Name(),
	// 	RawObject:  string(rawObject),
	// }, nil
}
