package mindre

import "github.com/Mind-Informatica-srl/mind-reminder/internal/logic"

type Reminder logic.Reminder

type Event logic.Event

type EventBase logic.EventBase

const (
	ActionCreate logic.Action = logic.ActionCreate
	ActionUpdate              = logic.ActionUpdate
	ActionDelete              = logic.ActionDelete
)
