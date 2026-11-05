package potager_test

import (
	"fmt"
	"time"

	"github.com/JMaxant/potager-app/internal/potager"
)

func ExamplePlanting_Status() {
	tomato := potager.Planting{
		Purpose:        potager.PurposeFood,
		MaturityFrom:   potager.MaturityFromTransplant,
		SownAt:         potager.DateFromInt(2026, time.March, 15),
		DaysToMaturity: 80,
	}

	// 1 April: the seedling is still in its tray.
	fmt.Println("01/04", tomato.Status(potager.DateFromInt(2026, time.April, 1)))

	// May 10th: the crop is transplanted and growing.
	tomato.TransplantedAt = potager.DateFromInt(2026, time.May, 10)
	fmt.Println("01/07", tomato.Status(potager.DateFromInt(2026, time.July, 1)))

	// July 29th: the crop is still growing.
	fmt.Println("29/07", tomato.Status(potager.DateFromInt(2026, time.July, 29)))

	// August 5th: the crop is due.
	fmt.Println("05/08", tomato.Status(potager.DateFromInt(2026, time.August, 5)))

	// August 15th: the crop is harvested.
	tomato.EndedAt = potager.DateFromInt(2026, time.August, 15)
	fmt.Println("15/08", tomato.Status(potager.DateFromInt(2026, time.August, 15)))

	// Output:
	// 01/04 planned
	// 01/07 growing
	// 29/07 growing
	// 05/08 due
	// 15/08 ended
}
