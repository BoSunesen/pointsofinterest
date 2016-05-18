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
	OpeningHours *OpeningHours
}

//TODO Use map, if it marshals to correct JSON, OpeningHours shouldn't be a pointer then
type OpeningHours struct {
	Sunday    []TimeInterval
	Monday    []TimeInterval
	Tuesday   []TimeInterval
	Wednesday []TimeInterval
	Thursday  []TimeInterval
	Friday    []TimeInterval
	Saturday  []TimeInterval
}

type TimeInterval struct {
	OpenFrom int
	OpenTo   int
}
