package metaweather_test

import (
	"strings"
	"testing"

	"github.com/themue/training-samples/pkg/metaweather"
)

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

func TestSuccessfulEmptyQuery(t *testing.T) {
	locations, err := metaweather.QueryLocations("thisissomestrangelocationwhichdoesnotexist")
	if err != nil {
		t.Fatalf("query returned error: %v", err)
	}
	if len(locations) != 0 {
		t.Fatalf("number of query results is %d", len(locations))
	}
}

func TestEmptyQueryValue(t *testing.T) {
	_, err := metaweather.QueryLocations("")
	if err == nil {
		t.Fatalf("error is nil")
	}
	if !strings.Contains(err.Error(), "JSON") {
		t.Fatalf("error is %v", err)
	}
}

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
}
