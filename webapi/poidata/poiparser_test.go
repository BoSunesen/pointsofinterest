package poidata

import (
	"testing"
	"time"
)

func TestParseMultiWeekdayOpenings(t *testing.T) {
	weekdayOpenings, err := parseMultiWeekdayOpenings("Tu/We/Fr:12AM-4AM;Fr-Su:11AM-1PM/11PM-1AM;Sa:12PM-4PM;")
	if err != nil {
		t.Error(err)
		return
	}
	if len(weekdayOpenings) != 6 {
		t.Errorf("Incorrect number of weekdays: %v", weekdayOpenings)
		return
	}

	//Sunday 11PM-1AM results in monday opening 0 - 1
	validateWeekdayOpening(t, weekdayOpenings, time.Monday, []int{0}, []int{1})
	validateWeekdayOpening(t, weekdayOpenings, time.Tuesday, []int{0}, []int{4})
	validateWeekdayOpening(t, weekdayOpenings, time.Wednesday, []int{0}, []int{4})
	validateWeekdayOpening(t, weekdayOpenings, time.Thursday, []int{}, []int{})
	validateWeekdayOpening(t, weekdayOpenings, time.Friday, []int{0, 11, 23}, []int{4, 13, 1})
	validateWeekdayOpening(t, weekdayOpenings, time.Saturday, []int{0, 11, 23, 12}, []int{1, 13, 1, 16})
	validateWeekdayOpening(t, weekdayOpenings, time.Sunday, []int{0, 11, 23}, []int{1, 13, 1})
}

func validateWeekdayOpening(t *testing.T, weekdayOpenings map[string][]TimeInterval, day time.Weekday, from, to []int) {
	openings := weekdayOpenings[day.String()]
	if len(openings) != len(from) && len(openings) != len(to) {
		t.Errorf("Incorrect number of openings for %v: %v", day.String(), openings)
		return
	}
	for i := 0; i < len(openings); i++ {
		isMissing := true
		for j := 0; j < len(openings); j++ {
			if openings[i].OpenFrom == from[j] && openings[i].OpenTo == to[j] {
				isMissing = false
			}
		}
		if isMissing {
			t.Errorf("Incorrect openings for %v: %v", day.String(), openings)
			return
		}
	}
}
