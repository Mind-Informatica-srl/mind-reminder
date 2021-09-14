// Package mrmodel expose the model to use remind package with
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
