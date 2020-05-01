package util

import (
	"math"
	"time"
)

type CalendarDay struct {
	Day       int   `json:"day" bson:"day"`
	Dayofweek int   `json:"day_of_week" bson:"day_of_week"`
	Works     bool  `json:"works" bson:"works"`
	Timestamp int64 `json:"-" bson:"timestamp"`
}

var firstwork time.Time

func CalculateCalendar(start time.Time, workdays int, weekdays int) []CalendarDay {
	var now = time.Now()
	var first = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	for first.Weekday() != 1 {
		first = first.AddDate(0, 0, -1)
	}

	firstwork = start

	var i = 0
	var calendr []CalendarDay
	for i < 42 {
		var day CalendarDay
		var date = first.AddDate(0, 0, i)
		day.Day = date.Day()
		day.Dayofweek = (int)(date.Weekday())
		day.Timestamp = date.Unix()
		day.Works = workday(date, workdays, weekdays)
		calendr = append(calendr, day)
		i++
	}
	return calendr
}

func workday(day time.Time, workdays int, weekdays int) bool {
	var span = day.Sub(firstwork).Hours() / 24
	if span < 0 {
		if span <= -(float64)(workdays+weekdays+1) {
			span = math.Abs(span) + (float64)(workdays-1)
		} else {
			span = span + (float64)(workdays+weekdays)
		}
	}
	return int(span)%(weekdays+workdays) <= workdays-1
}
