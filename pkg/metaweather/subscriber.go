package metaweather

import (
	"context"
	"log"
	"strings"
	"time"
)

// Subscriber manages the MetaWeather subscriptions and
// chronologically polls the data.
type Subscriber struct {
	ctx       context.Context
	interval  time.Duration
	actionc   chan func()
	locations map[string]int
	weathers  map[int]Weather
}

// StartSubscriber makes the Subscriber run in the background.
func StartSubscriber(ctx context.Context, interval time.Duration) *Subscriber {
	s := &Subscriber{
		ctx:       ctx,
		interval:  interval,
		actionc:   make(chan func(), 16),
		locations: make(map[string]int),
		weathers:  make(map[int]Weather),
	}
	go s.backend()
	return s
}

// Subscribe adds the subscription of one or multiple cities. Their
// names will be returned.
func (s *Subscriber) Subscribe(query string) []string {
	names := []string{}

	s.doSync(func() {
		locations, err := QueryLocations(query)
		if err != nil {
			log.Printf("query of %q failed: %v", query, err)
			return
		}
		for _, location := range locations {
			name := strings.ToLower(location.Title)
			s.locations[name] = location.WOEID
			names = append(names, name)
			if _, ok := s.weathers[location.WOEID]; !ok {
				// It's new, so add it.
				weather, err := ReadWeather(location.WOEID)
				if err != nil {
					log.Printf("subscription of %q failed: %v", name, err)
					continue
				}
				s.weathers[location.WOEID] = weather
				log.Printf("subscribed %q", name)
			}
		}
	})

	return names
}

// Fetch retrieves a number of Weathers. Any so far unsubscribed name
// will be ignored.
func (s *Subscriber) Fetch(names ...string) []Weather {
	weathers := []Weather{}

	s.doSync(func() {
		for _, name := range names {
			woeid, ok := s.locations[strings.ToLower(name)]
			if !ok {
				continue
			}
			weathers = append(weathers, s.weathers[woeid])
		}
	})

	return weathers
}

// doSync sends an action for execution to the backend and waits
// until its done. Right now no handling of timeouts to cover
// closed channels.
func (s *Subscriber) doSync(action func()) {
	donec := make(chan struct{})

	s.actionc <- func() {
		action()
		close(donec)
	}

	<-donec
}

// backend is the goroutine of the Subscriber.
func (s *Subscriber) backend() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	defer close(s.actionc)
	for {
		select {
		case <-s.ctx.Done():
			return
		case action := <-s.actionc:
			action()
		case <-ticker.C:
			s.updateAll()
		}
	}
}

// updateAll starts a goroutine per location that reads the weather
// information and stores it. This way the other interaction with
// the subscriber is not blocked.
func (s *Subscriber) updateAll() {
	for woeid := range s.weathers {
		go s.updateOne(woeid)
	}
}

// updateOne updates the weather data for one location.
func (s *Subscriber) updateOne(woeid int) {
	s.actionc <- func() {
		name := s.weathers[woeid].Title
		log.Printf("updating weather of %q...", name)
		weather, err := ReadWeather(woeid)
		if err != nil {
			// Don't care a lot, just log it.
			log.Printf("updating weather of %q failed: %v", name, err)
			return
		}
		s.weathers[woeid] = weather
	}
}
