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
	Applicant      string
	Address        string
	Dayshours      string
	Facilitytype   string
	Fooditems      string
	Latitude       string
	Longitude      string
	Status         string
	LatitudeFloat  float64
	LongitudeFloat float64
	OpeningHours   *WeekdayOpenings
}

type WeekdayOpenings struct {
	Sunday    []OpeningHours
	Monday    []OpeningHours
	Tuesday   []OpeningHours
	Wednesday []OpeningHours
	Thursday  []OpeningHours
	Friday    []OpeningHours
	Saturday  []OpeningHours
}

type OpeningHours struct {
	OpenFrom int
	OpenTo   int
}
