package poidata

import "time"

//TODO Test filters
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
	weekdayString := time.Weekday(weekdayInt).String()
	for _, element := range *parsedData {
		includeElement := false
		if element.OpeningHours != nil {
			for _, interval := range element.OpeningHours[weekdayString] {
				if interval.OpenFrom <= hour && interval.OpenTo > hour {
					includeElement = true
					break
				}
			}
		}
		if includeElement {
			openingsFilteredData = append(openingsFilteredData, element)
		}
	}
	return &openingsFilteredData
}
