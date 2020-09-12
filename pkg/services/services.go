package services

import (
	"context"
	"log"
)

// Service defines a service component a User can book
// and execute.
type Service interface {
	// ID returns a helping identifier of the Service.
	ID() string

	// Do executes the Service.
	Do() error
}

// Services is a collection of services able to be spawned.
type Services map[string]Service

// Spawn executes all services in the background.
func (svcs Services) Spawn() {
	go func() {
		for id, svc := range svcs {
			err := svc.Do()
			if err != nil {
				log.Printf("execution of service %q failed: %v", id, err)
			}
		}
	}()
}

// ServiceProvider manages the Services per consumer. Those
// can be added and removed as well as spawned. In that
// case the individual services are executed concurrently.
type ServiceProvider struct {
	ctx      context.Context
	actionc  chan func()
	bookings map[string]Services
}

// StartServiceProvider makes the ServiceProvider run in the background.
func StartServiceProvider(ctx context.Context) *ServiceProvider {
	sp := &ServiceProvider{
		ctx:      ctx,
		actionc:  make(chan func(), 16),
		bookings: make(map[string]Services),
	}
	go sp.backend()
	return sp
}

// Book assigns Services to a consumer, like e.g. a User.
func (sp *ServiceProvider) Book(id string, svcs ...Service) {
	sp.doSync(func() {
		current, ok := sp.bookings[id]
		if !ok {
			current = make(Services)
		}
		for _, svc := range svcs {
			current[svc.ID()] = svc
		}
		sp.bookings[id] = current
	})
}

// Unbook drops assignment of a Services to a consumer.
func (sp *ServiceProvider) Unbook(id string, svcIDs ...string) {
	sp.doSync(func() {
		current, ok := sp.bookings[id]
		if !ok {
			return
		}
		for _, svcID := range svcIDs {
			delete(current, svcID)
		}
		if len(current) == 0 {
			delete(sp.bookings, id)
		} else {
			sp.bookings[id] = current
		}
	})
}

// doSync sends an action for execution to the backend and waits
// until its done. Right now no handling of timeouts to cover.
// closed channels.
func (sp *ServiceProvider) doSync(action func()) {
	donec := make(chan struct{})

	sp.actionc <- func() {
		action()
		close(donec)
	}

	<-donec
}

// backend is the goroutine of the ServiceProvider.
func (sp *ServiceProvider) backend() {
	defer close(sp.actionc)
	for {
		select {
		case <-sp.ctx.Done():
			return
		case action := <-sp.actionc:
			action()
		}
	}
}
