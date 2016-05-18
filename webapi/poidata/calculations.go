package poidata

import "math"

type BoundingBox struct {
	MaxLatitude  float64
	MinLatitude  float64
	MaxLongitude float64
	MinLongitude float64
}

func CalculateBoundingBox(latitude float64, longitude float64, distance int) BoundingBox {
	const earthRadius float64 = 6371000
	distanceFloat := float64(distance)

	latitudeRadians := ConvertDegreesToRadians(latitude)
	latitudeDiff := distanceFloat / earthRadius

	longitudeRadians := ConvertDegreesToRadians(longitude)
	longitudeDiff := distanceFloat / (math.Cos(latitudeRadians) * earthRadius)

	maxLatitude := latitudeRadians + latitudeDiff
	minLatitude := latitudeRadians - latitudeDiff
	maxLongitude := longitudeRadians + longitudeDiff
	minLongitude := longitudeRadians - longitudeDiff
	return BoundingBox{
		MaxLatitude:  ConvertRadiansToDegrees(maxLatitude),
		MinLatitude:  ConvertRadiansToDegrees(minLatitude),
		MaxLongitude: ConvertRadiansToDegrees(maxLongitude),
		MinLongitude: ConvertRadiansToDegrees(minLongitude),
	}
}

func ConvertDegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func ConvertRadiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}
