package potager

import (
	"time"
)

var locParis *time.Location

func init() {
	var err error
	locParis, err = time.LoadLocation("Europe/Paris")
	if err != nil {
		panic(err)
	}
}

// Date represents a date without time information.
type Date struct {
	t time.Time
}

func (d Date) String() string {
	return d.t.Format(time.DateOnly)
}

// DateFromTime returns a Date in Europe/Paris timezone created from a time.Time object.
func DateFromTime(date time.Time) Date {
	date = date.In(locParis)
	year, month, day := date.Date()
	return DateFromInt(year, month, day)
}

// DateFromInt returns a Date created from a year, month and day.
func DateFromInt(year int, month time.Month, day int) Date {
	// Timezone is not needed here, it is only relevant when extracting date from a Time object to normalize them.
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func (d Date) addDays(days int) Date {
	return Date{d.t.AddDate(0, 0, days)}
}

//func (d Date) compare(date Date) bool {
//	return d.t.Equal(date.t)
//}
