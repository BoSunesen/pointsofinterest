package poidata

import "time"

func FilterByLocation(parsedData *[]ParsedPoiData, latitude, longitude float64, distance int) *[]ParsedPoiData {
	geoFilteredData := make([]ParsedPoiData, 0, len(*parsedData))
	boundingBox := CalculateBoundingBox(latitude, longitude, distance)
	for _, v := range *parsedData {
		if v.Latitude <= boundingBox.MaxLatitude &&
			v.Latitude >= boundingBox.MinLatitude &&
			v.Longitude <= boundingBox.MaxLongitude &&
			v.Longitude >= boundingBox.MinLongitude {
			geoFilteredData = append(geoFilteredData, v)
		}
	}
	return &geoFilteredData
}

func FilterByOpeningHours(parsedData *[]ParsedPoiData, weekdayInt, hour int) *[]ParsedPoiData {
	openingsFilteredData := make([]ParsedPoiData, 0, len(*parsedData))
	var weekdays []string
	if weekdayInt >= 0 && weekdayInt < 7 {
		weekdays = []string{time.Weekday(weekdayInt).String()}
	} else {
		weekdays = make([]string, 7)
		for i := 0; i < 7; i++ {
			weekdays[i] = time.Weekday(i).String()
		}
	}

	for _, element := range *parsedData {
		includeElement := false
		if element.OpeningHours != nil {
			for _, weekdayString := range weekdays {
				for _, interval := range element.OpeningHours[weekdayString] {
					if isHourInInterval(hour, interval) {
						includeElement = true
						break
					}
				}
			}
		}
		if includeElement {
			openingsFilteredData = append(openingsFilteredData, element)
		}
	}
	return &openingsFilteredData
}

func isHourInInterval(hour int, interval TimeInterval) bool {
	if hour < 0 || hour >= 24 {
		//Accept any interval if hour is outside legal values
		return true
	}
	if interval.OpenFrom > hour {
		//Not open yet
		return false
	}
	if interval.OpenFrom > interval.OpenTo {
		//Open until tomorrow
		return true
	}
	if interval.OpenTo > hour {
		//Still open
		return true
	}
	//Closed
	return false
}
