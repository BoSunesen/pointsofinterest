package poidata

import "testing"

func TestFilterByLocation(t *testing.T) {
	poiData := []ParsedPoiData{
		ParsedPoiData{
			Applicant: "Test1a",
			Latitude:  10.09,
			Longitude: 11.11,
		},
		ParsedPoiData{
			Applicant: "Test1b",
			Latitude:  10.11,
			Longitude: 11.09,
		},
		ParsedPoiData{
			Applicant: "Test2",
			Latitude:  -12.3,
			Longitude: -13.4,
		},
		ParsedPoiData{
			Applicant: "Test3",
			Latitude:  -14.5,
			Longitude: 15.6,
		},
		ParsedPoiData{
			Applicant: "Test4",
			Latitude:  16.7,
			Longitude: -17.8,
		},
	}
	filtered := FilterByLocation(&poiData, 10.1, 11.1, 10000)
	numberOfLocations := len(*filtered)
	if numberOfLocations != 2 {
		t.Errorf("Incorrect number of locations: %v", numberOfLocations)
	}
	for _, v := range *filtered {
		if v.Applicant != "Test1a" && v.Applicant != "Test1b" {
			t.Errorf("Incorrect locations: %v", v.Applicant)
		}
	}
}

func TestFilterByOpeningHours(t *testing.T) {
	poiData := []ParsedPoiData{
		ParsedPoiData{
			Applicant:    "Test1",
			OpeningHours: nil,
		},
		ParsedPoiData{
			Applicant: "Test2",
			OpeningHours: map[string][]TimeInterval{
				"Monday": []TimeInterval{TimeInterval{OpenFrom: 0, OpenTo: 23}},
			},
		},
		ParsedPoiData{
			Applicant: "Test3",
			OpeningHours: map[string][]TimeInterval{
				"Friday": []TimeInterval{TimeInterval{OpenFrom: 22, OpenTo: 23}},
			},
		},
		ParsedPoiData{
			Applicant: "Test4",
			OpeningHours: map[string][]TimeInterval{
				"Friday": []TimeInterval{TimeInterval{OpenFrom: 21, OpenTo: 2}},
			},
		},
		ParsedPoiData{
			Applicant: "Test5",
			OpeningHours: map[string][]TimeInterval{
				"Friday": []TimeInterval{TimeInterval{OpenFrom: 18, OpenTo: 22}},
			},
		},
		ParsedPoiData{
			Applicant: "Test6",
			OpeningHours: map[string][]TimeInterval{
				"Friday": []TimeInterval{TimeInterval{OpenFrom: 9, OpenTo: 19}},
			},
		},
	}
	filtered := FilterByOpeningHours(&poiData, 5, 21)
	numberOfLocations := len(*filtered)
	if numberOfLocations != 2 {
		t.Errorf("Incorrect number of locations: %v", numberOfLocations)
	}
	for _, v := range *filtered {
		if v.Applicant != "Test4" && v.Applicant != "Test5" {
			t.Errorf("Incorrect locations: %v", v.Applicant)
		}
	}

}
