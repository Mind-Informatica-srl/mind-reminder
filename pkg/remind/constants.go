// Package remind expose the model to use remind package with
package remind

// Action is the type representing what you can do with a record: create, update, delete
type Action string

const (
	// ActionCreate represent the creating action
	ActionCreate Action = "create"
	// ActionUpdate represent the updating action
	ActionUpdate Action = "update"
	// ActionDelete represent the deleting action
	ActionDelete Action = "delete"
)

type DataFieldType string

const (
	Text     DataFieldType = "text"
	Number   DataFieldType = "number"
	Date     DataFieldType = "date"
	Select   DataFieldType = "select"
	Radio    DataFieldType = "radio"
	Checkbox DataFieldType = "checkbox"
)

type SectionConfiguration string

const (
	EventFirst  SectionConfiguration = "event_first"
	ObjectFirst SectionConfiguration = "object_first"
	// Both SectionConfiguration = "both"
)
