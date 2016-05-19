package poidata

import (
	"errors"
	"fmt"
	"github.com/BoSunesen/pointsofinterest/webapi/logging"
	"golang.org/x/net/context"
	"strconv"
	"strings"
	"time"
)

type PoiParser struct {
	Logger logging.Logger
}

//TODO Test ParsePoiData
//"Applicant": "Mike's Catering"
//"Address": "860 BROADWAY"
//"Dayshours": "Mo/Tu/We/Th/Fr:7AM-8AM;Mo/Mo/Tu/Tu/We:9AM-11AM;Su:9AM-2PM;Sa:9AM-3PM;Mo/Mo/Tu/Tu/We:11AM-1PM;Mo-Fr:"
func (parser PoiParser) ParsePoiData(ctx context.Context, input *[]PoiData) *[]ParsedPoiData {
	output := make([]ParsedPoiData, len(*input))
	skipped := 0
	for i, poi := range *input {
		latitude, err := strconv.ParseFloat(poi.Latitude, 64)
		if err != nil {
			if len(poi.Longitude) > 1 {
				parser.Logger.Errorf(ctx, "Error while parsing latitude \"%v\": %v", poi.Latitude, err)
			} else {
				//TODO parser.Logger.Infof(ctx, "Empty latitude")
			}
			skipped++
			continue
		}
		longitude, err := strconv.ParseFloat(poi.Longitude, 64)
		if err != nil {
			if len(poi.Longitude) > 1 {
				parser.Logger.Errorf(ctx, "Error while parsing longitude \"%v\": %v", poi.Longitude, err)
			} else {
				//TODO parser.Logger.Infof(ctx, "Empty longitude")
			}
			skipped++
			continue
		}

		weekdayOpenings, err := parser.parseMultiWeekdayOpenings(poi.Dayshours)
		if err != nil {
			//TODO parser.Logger.Errorf(ctx, "Error while parsing opening hours \"%v\": %v", poi.Dayshours, err)
			weekdayOpenings = make(map[string][]TimeInterval, 0)
		}

		parsed := ParsedPoiData{
			Applicant:    poi.Applicant,
			Address:      poi.Address,
			Dayshours:    poi.Dayshours,
			Facilitytype: poi.Facilitytype,
			Fooditems:    poi.Fooditems,
			Status:       poi.Status,
			Latitude:     latitude,
			Longitude:    longitude,
			OpeningHours: weekdayOpenings,
		}
		output[i-skipped] = parsed
	}
	result := output[:len(output)-skipped]
	return &result
}

func (parser PoiParser) parseMultiWeekdayOpenings(weekdayOpeningsString string) (map[string][]TimeInterval, error) {
	weekdayOpenings := make(map[string][]TimeInterval, 7)
	weekdayOpeningsSplit := strings.Split(weekdayOpeningsString, ";")
	for _, v := range weekdayOpeningsSplit {
		err := parser.parseSingleWeekdayOpenings(v, weekdayOpenings)
		if err != nil {
			return nil, err
		}
	}
	return weekdayOpenings, nil
}

func (parser PoiParser) parseSingleWeekdayOpenings(weekdayOpeningsString string, weekdayOpenings map[string][]TimeInterval) error {
	days := make(map[time.Weekday]bool, 7)

	weekdaysAndTimes := strings.Split(weekdayOpeningsString, ":")
	if len(weekdaysAndTimes) != 2 {
		return errors.New("Weekday openings did not contain exactly one ':'")
	}

	currentString := weekdaysAndTimes[0]
	for len(currentString) > 0 {
		if len(currentString) < 2 {
			return fmt.Errorf("Rest of string (%v) is not long enough to contain day", currentString)
		}
		var weekdayString string
		weekdayString, currentString = currentString[:2], currentString[2:]

		weekday, err := parseWeekday(weekdayString)
		if err != nil {
			return err
		}

		days[weekday] = true

		if len(currentString) > 0 {
			var weekdaySplitter string
			weekdaySplitter, currentString = currentString[:1], currentString[1:]

			switch weekdaySplitter {
			case "/":
				continue
			case "-":
				if weekday == time.Sunday {
					return errors.New("Days interval starting on sunday is not supported")
				}

				if len(currentString) < 2 {
					return fmt.Errorf("Rest of string (%v) is not long enough to contain interval end day", currentString)
				}
				var intervalEndWeekdayString string
				intervalEndWeekdayString, currentString = currentString[:2], currentString[2:]
				weekdayEnd, err := parseWeekday(intervalEndWeekdayString)
				if err != nil {
					return err
				}

				var end int
				if weekdayEnd == time.Sunday {
					end = 7
				} else {
					end = int(weekdayEnd)
				}
				for i := int(weekday); i <= end; i++ {
					var nextDay time.Weekday
					if i == 7 {
						nextDay = time.Sunday
					} else {
						nextDay = time.Weekday(i)
					}
					days[nextDay] = true
				}
			default:
				return fmt.Errorf("Illegal weekday splitter: %v", weekdaySplitter)
			}
		}
	}
	err := parser.addTimeSlot(weekdaysAndTimes[1], days, weekdayOpenings)
	if err != nil {
		return err
	}
	return nil
}

func (parser PoiParser) addTimeSlot(timeString string, days map[time.Weekday]bool, weekdayOpenings map[string][]TimeInterval) error {
	openings, err := parser.parseTime(timeString)
	if err != nil {
		return err
	}

	for weekday, dayIsSpecified := range days {
		if dayIsSpecified {
			weekdayOpenings[weekday.String()] = append(weekdayOpenings[weekday.String()], openings...)
		}
	}

	return nil
}

func (parser PoiParser) parseTime(timeString string) ([]TimeInterval, error) {
	openingHours := make([]TimeInterval, 0)
	if len(timeString) == 0 {
		return nil, errors.New("Time string was empty")
	}

	timeStringElements := strings.Split(timeString, "/")
	for _, timeStringElement := range timeStringElements {
		fromAndToStrings := strings.Split(timeStringElement, "-")
		if len(fromAndToStrings) != 2 {
			return nil, fmt.Errorf("Time interval (%v) does not contain exactly one '-'", timeStringElement)
		}
		from, err := parseTimeOfDay(fromAndToStrings[0])
		if err != nil {
			return nil, err
		}
		to, err := parseTimeOfDay(fromAndToStrings[1])
		if err != nil {
			return nil, err
		}
		openingHours = append(openingHours, TimeInterval{from, to})
	}
	return openingHours, nil
}

func parseTimeOfDay(timeOfDayString string) (int, error) {
	fromTimeStringSplitIndex := len(timeOfDayString) - 2
	if fromTimeStringSplitIndex < 1 {
		return 0, fmt.Errorf("Time of day is too short: %v", timeOfDayString)
	}
	hourString, amOrPm := timeOfDayString[:fromTimeStringSplitIndex], timeOfDayString[fromTimeStringSplitIndex:]
	hour, err := strconv.Atoi(hourString)
	if err != nil {
		return 0, err
	}
	if amOrPm == "AM" {
		if hour == 12 {
			hour = 0
		}
	} else if amOrPm == "PM" {
		if hour != 12 {
			hour = hour + 12
		}
	} else {
		return 0, fmt.Errorf("Expected 'AM' or 'PM' but found '%v'", amOrPm)
	}
	return hour, nil
}

func parseWeekday(weekday string) (time.Weekday, error) {
	switch weekday {
	case "Su":
		return time.Sunday, nil
	case "Mo":
		return time.Monday, nil
	case "Tu":
		return time.Tuesday, nil
	case "We":
		return time.Wednesday, nil
	case "Th":
		return time.Thursday, nil
	case "Fr":
		return time.Friday, nil
	case "Sa":
		return time.Saturday, nil
	default:
		return time.Sunday, fmt.Errorf("Unknown weekday string: %v", weekday)
	}
}