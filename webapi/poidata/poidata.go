package poidata

type PoiData struct {
	Applicant    string
	Address      string
	Dayshours    string
	Facilitytype string
	Fooditems    string
	Latitude     string
	Longitude    string
	Status       string
}

type ParsedPoiData struct {
	Applicant    string
	Address      string
	Dayshours    string
	Facilitytype string
	Fooditems    string
	Status       string
	Latitude     float64
	Longitude    float64
	OpeningHours map[string][]TimeInterval
}

type TimeInterval struct {
	OpenFrom int
	OpenTo   int
}
