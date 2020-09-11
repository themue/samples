package metaweather

import (
	"time"
)

// Location is part of the locations response.
type Location struct {
	Title        string `json:"title"`
	LocationType string `json:"location_type"`
	LattLong     string `json:"latt_long"`
	WOEID        int    `json:"woeid"`
	Distance     int    `json:"distance"`
}

// Locations is returned when querying for locations.
type Locations []Location

// Parent is part of the location resonse and describes the
// parent location."
type Parent struct {
	Title        string `json:"title"`
	LocationType string `json:"location_type"`
	LattLong     string `json:"latt_long"`
	WOEID        int    `json:"woeid"`
}

// ConsolidatedWeather is part of the location response and
// contains the consolidated weather information.
type ConsolidatedWeather struct {
	ID                   int     `json:"id"`
	ApplicableDate       string  `json:"applicable_date"`
	WeatherStateName     string  `json:"weather_state_name"`
	WeatherStateAbbr     string  `json:"weather_state_abbr"`
	WindSpeed            float64 `json:"wind_speed"`
	WindDirection        float64 `json:"wind_direction"`
	WindDirectionCompass string  `json:"wind_direction_compass"`
	MinTemp              float64 `json:"min_temp"`
	MaxTemp              float64 `json:"max_temp"`
	TheTemp              float64 `json:"the_temp"`
	AirPresure           float64 `json:"air_presure"`
	Humidity             float64 `json:"humidity"`
	Visibility           float64 `json:"visibility"`
	Predictability       int     `json:"predictability"`
}

// Source is part of the location response and here one of the
// returned sources of the weather information.
type Source struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Weather is returned when querying a woeid.
type Weather struct {
	Title               string                `json:"title"`
	LocationType        string                `json:"location_type"`
	LattLong            string                `json:"latt_long"`
	Time                time.Time             `json:"time"`
	Sunrise             time.Time             `json:"sun_rise"`
	Sunset              time.Time             `json:"sun_set"`
	TimezoneName        string                `json:"timezone_name"`
	Parent              Parent                `json:"parent"`
	ConsolidatedWeather []ConsolidatedWeather `json:"consolidated_weather"`
	Sources             []Source              `json:"sources"`
}
