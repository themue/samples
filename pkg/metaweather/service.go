package metaweather

// Callback defines what has to be passed to a new Service
// to handle retrieved Weathers.
type Callback func([]Weather)

// Service implements services.Service for MetaWeather. A consumer
// like a User can use it to book the subscription to one or more
// locations.
type Service struct {
	id       string
	sub      *Subscriber
	callback Callback
	names    []string
}

// NewService creates a MetaWeather service instance.
func NewService(
	id string,
	sub *Subscriber,
	callback Callback,
	names ...string,
) *Service {
	return &Service{
		id:       id,
		sub:      sub,
		callback: callback,
		names:    names,
	}
}

// ID implements services.Service.
func (s *Service) ID() string {
	return s.id
}

// Do implements services.Service.
func (s *Service) Do() error {
	ws := s.sub.Fetch(s.names...)
	s.callback(ws)
}
