package potager

import (
	"time"
)

type Date struct {
	t time.Time
}

func DateFromTime(date time.Time) Date {
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		panic(err)
	}
	date = date.In(loc)
	year, month, day := date.Date()
	return DateFromInt(year, month, day)
}

func DateFromInt(year int, month time.Month, day int) Date {
	// Timezone is not needed here, it is only relevant when extracting date from a Time object to normalize them.
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

//func (d Date) addDays(days int) Date {
//	panic("todo")
//}

//func (d Date) compare(date Date) bool {
//	return d.t.Equal(date.t)
//}
