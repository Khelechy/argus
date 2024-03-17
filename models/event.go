package models

import (
	"time"

	"github.com/khelechy/argus/enums"
)

type Event struct {
	Action enums.Action
	ActionDescription string
	Name string
	Timestamp time.Time
}
