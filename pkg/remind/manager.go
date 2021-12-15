package remind

import "gorm.io/gorm"

type manager struct {
	eventTableName        string
	remindTableName       string
	accomplisherTableName string
}

// AddEvent aggiunge un evento
func (m *manager) AddEvent(db *gorm.DB, event *Event) (err error) {
	event.manager = m
	err = db.Create(event).Error
	return
}

// UpdateEvent modifica un evento
func (m *manager) UpdateEvent(db *gorm.DB, event *Event) (err error) {
	event.manager = m
	err = db.Save(event).Error
	return
}

// DeleteEvent elimina un evento
func (m *manager) DeleteEvent(db *gorm.DB, event *Event) (err error) {
	event.manager = m
	err = db.Delete(event).Error
	return
}

func CreateManager(eventTableName string, remindTableName string, accomplisherTableName string) manager {
	return manager{
		eventTableName:        eventTableName,
		remindTableName:       remindTableName,
		accomplisherTableName: accomplisherTableName,
	}
}
