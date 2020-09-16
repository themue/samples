package metaweather_test

import (
	"strings"
	"testing"

	"github.com/themue/training-samples/pkg/metaweather"
)

// TestSuccessfulSingleQuery verifies the query of a location
// with only one matching result. In case of an error at
// QueryLocations it's due to some kind of network trouble.
func TestSuccessfulSingleQuery(t *testing.T) {
	locations, err := metaweather.QueryLocations("london")
	if err != nil {
		t.Fatalf("query returned error: %v", err)
	}
	if len(locations) != 1 {
		t.Fatalf("number of query results is %d", len(locations))
	}
	if locations[0].Title != "London" {
		t.Fatalf("query results returned wrong location: %v", locations[0])
	}
}

// TestSuccessfulMultiQuery verifies the query of multiple
// locations based on parts of it.
func TestSuccessfulMultiQuery(t *testing.T) {
	locations, err := metaweather.QueryLocations("san")
	if err != nil {
		t.Fatalf("query returned error: %v", err)
	}
	if !(len(locations) > 1) {
		t.Fatalf("number of query results is %d", len(locations))
	}
	for i, loc := range locations {
		if !strings.Contains(strings.ToLower(loc.Title), "san") {
			t.Fatalf("query result %d contains no 'san': %v", i, loc)
		}
	}
}

// TestSuccessfulEmptyQuery verifies the query with no locations
// but also no error as result.
func TestSuccessfulEmptyQuery(t *testing.T) {
	locations, err := metaweather.QueryLocations("thisissomestrangelocationwhichdoesnotexist")
	if err != nil {
		t.Fatalf("query returned error: %v", err)
	}
	if len(locations) != 0 {
		t.Fatalf("number of query results is %d", len(locations))
	}
}

// TestEmptyQueryValue verifies the invalid usage of an
// empty query string.
func TestEmptyQueryValue(t *testing.T) {
	_, err := metaweather.QueryLocations("")
	if err == nil {
		t.Fatalf("error is nil")
	}
	if !strings.Contains(err.Error(), "JSON") {
		t.Fatalf("error is %v", err)
	}
}

// TestReadLondon verifies the retrieval of the weather data
// of one location by its "Where On Earth ID" (woeid).
func TestReadLondon(t *testing.T) {
	locations, err := metaweather.QueryLocations("london")
	if err != nil {
		t.Fatalf("query returned error: %v", err)
	}
	woeid := locations[0].WOEID
	london, err := metaweather.ReadWeather(woeid)
	if err != nil {
		t.Fatalf("read returned error: %v", err)
	}
	minTemp := london.ConsolidatedWeather[0].MinTemp
	maxTemp := london.ConsolidatedWeather[0].MaxTemp
	theTemp := london.ConsolidatedWeather[0].TheTemp
	t.Logf("London Temperature: Min = %.2f / Max = %.2f / Act = %.2f", minTemp, maxTemp, theTemp)
	if theTemp < minTemp || theTemp > maxTemp {
		t.Fatalf("temparature values are invalid ???")
	}
}
