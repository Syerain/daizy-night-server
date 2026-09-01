package v1

import "time"

// CalendarGetResponse renders a timetable without leaking internal ids or
// soft-delete bookkeeping (explicit mapping, same idea as InfoMe).
type CalendarGetResponse struct {
	CalendarID uint                   `json:"calendar_id"`
	Records    []CalendarItemResponse `json:"records"`
}

type CalendarItemResponse struct {
	Weekday  time.Weekday `json:"weekday"`
	StartMin int          `json:"start_min"`
	EndMin   int          `json:"end_min"`
	Title    string       `json:"title"`
}

type CalendarPutResponse struct {
	Message string `json:"message"`
}

type CalendarDeleteResponse struct {
	Message string `json:"message"`
}
