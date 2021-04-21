package mindreminder

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

const (
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"
)

type UpdateDiff map[string]interface{}

//callback da eseguire dopo after_create
func (p *Plugin) addCreated(scope *gorm.Scope) {
	if isLoggable(scope.Value) && isEnabled(scope.Value) {
		_ = addRecord(scope, actionCreate)
	}
}

//callback da eseguire dopo after_update
func (p *Plugin) addUpdated(scope *gorm.Scope) {
	if !isLoggable(scope.Value) || !isEnabled(scope.Value) {
		return
	}

	// if p.opts.lazyUpdate {
	// 	record, err := p.GetLastRecord(interfaceToString(scope.PrimaryKeyValue()), false)
	// 	if err == nil {
	// 		if isEqual(record.RawObject, scope.Value, p.opts.lazyUpdateFields...) {
	// 			return
	// 		}
	// 	}
	// }

	_ = addUpdateRecord(scope, p.opts)
}

//callback da eseguire dopo after_delete.
func (p *Plugin) addDeleted(scope *gorm.Scope) {
	if isLoggable(scope.Value) && isEnabled(scope.Value) {
		_ = addRecord(scope, actionDelete)
	}
}

// Writes new reminder row to db.
func addRecord(scope *gorm.Scope, action string) error {
	cl, err := newReminder(scope, action)
	if err != nil {
		return nil
	}

	return scope.DB().Create(cl).Error
}

func addUpdateRecord(scope *gorm.Scope, opts options) error {
	cl, err := newReminder(scope, actionUpdate)
	if err != nil {
		return err
	}

	return scope.DB().Create(cl).Error
}

func newReminder(scope *gorm.Scope, action string) (*Reminder, error) {
	rawObject, err := json.Marshal(scope.Value)
	if err != nil {
		return nil, err
	}
	dateReminder, err := generateReminder(scope)
	if err != nil {
		return nil, err
	}
	return &Reminder{
		Action:     action,
		ObjectID:   interfaceToString(scope.PrimaryKeyValue()),
		ObjectType: scope.GetModelStruct().ModelType.Name(),
		RawObject:  string(rawObject),
		RawMeta:    string(fetchChangeLogMeta(scope)),
		RemindAt:   *dateReminder,
	}, nil
}
