package response

import (
	"github.com/fufuceng/gexc/time"
	"github.com/fufuceng/gexc/types"
)

//History is representation of the
//response of the history function result
type History struct {
	Base    string             `json:"base"`
	StartAt time.Gexc          `json:"start_at"`
	EndAt   time.Gexc          `json:"end_at"`
	Rates   types.TimeRateItem `json:"rates"`
}

//SingleDate is representation of the
//response of the latest function result
type SingleDate struct {
	Base  string         `json:"base"`
	Rates types.RateItem `json:"rates"`
	Date  time.Gexc      `json:"date"`
}
