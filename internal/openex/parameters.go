package openex

import "github.com/fufuceng/gexc/time"

type LatestParams struct {
	Base    string   `json:"base" url:"base"`
	Symbols []string `json:"symbols" url:"symbols"`
}

type SingleDateParams struct {
	Date    time.Gexc `json:"date"`
	Base    string    `json:"base" url:"base"`
	Symbols []string  `json:"symbols" url:"symbols"`
}

type HistoryParams struct {
	StartAt time.Gexc `json:"start_at" url:"start_at"`
	EndAt   time.Gexc `json:"end_at" url:"end_at"`
	Base    string    `json:"base" url:"base"`
	Symbols []string  `json:"symbols" url:"symbols"`
}
