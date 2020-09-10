package metaweather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	queryURL = "https://www.metaweather.com/api/location/search/?query=%s"
	readURL  = "https://www.metaweather.com/api/location/%d/"
)

// QueryLocations queries MetaWeather for locations with matching titles
// or title parts.
func QueryLocations(query string) (Locations, error) {
	url := fmt.Sprintf(queryURL, query)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot send MetaWeather query: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve body: %v", err)
	}

	var locations Locations

	err = json.Unmarshal(body, &locations)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %v", err)
	}

	return locations, nil
}

// ReadWeather retrieves the location and weather information for the
// given Where On Earth ID.
func ReadWeather(woeid int) (Weather, error) {
	var weather Weather

	url := fmt.Sprintf(readURL, woeid)
	resp, err := http.Get(url)
	if err != nil {
		return weather, fmt.Errorf("cannot read weather: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return weather, fmt.Errorf("cannot retrieve body: %v", err)
	}

	err = json.Unmarshal(body, &weather)
	if err != nil {
		return weather, fmt.Errorf("cannot unmarshal response: %v", err)
	}

	return weather, nil
}
