package mindre

import "github.com/Mind-Informatica-srl/mind-reminder/internal/logic"

type Reminder logic.Reminder

type Event logic.Event

type EventBase logic.EventBase

type Action logic.Action

const (
	ActionCreate Action = Action(logic.ActionCreate)
	ActionUpdate        = Action(logic.ActionUpdate)
	ActionDelete        = Action(logic.ActionDelete)
)
