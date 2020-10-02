// The Samples Project
//
// Copyright 2020 Frank Mueller / Oldenburg / Germany / World
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.
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

// Spawn executes all services concurrently.
func (svcs Services) Spawn() {
	go func() {
		for id, svc := range svcs {
			// Don't use loop variables directly, they will
			// change during iteration.
			go func(doID string, do func() error) {
				if err := do(); err != nil {
					log.Printf("execution of service %q failed: %v", doID, err)
				}
			}(id, svc.Do)
		}
	}()
}

// Provider manages the Services per consumer. Those
// can be added and removed as well as spawned. In that
// case the individual services are executed concurrently.
type Provider struct {
	ctx      context.Context
	actionc  chan func()
	bookings map[string]Services
}

// StartProvider creates a Provider running as goroutine.
func StartProvider(ctx context.Context) *Provider {
	p := &Provider{
		ctx:      ctx,
		actionc:  make(chan func(), 16),
		bookings: make(map[string]Services),
	}
	go p.backend()
	return p
}

// Book assigns Services to a consumer, like e.g. a User. The
// number of booked services is returned.
func (p *Provider) Book(consumerID string, svcs ...Service) int {
	var svcCnt int
	p.doSync(func() {
		current, ok := p.bookings[consumerID]
		if !ok {
			current = make(Services)
		}
		for _, svc := range svcs {
			current[svc.ID()] = svc
		}
		svcCnt = len(current)
		p.bookings[consumerID] = current
	})
	return svcCnt
}

// Unbook drops assignment of a Services to a consumer. The
// number of booked services is returned.
func (p *Provider) Unbook(consumerID string, svcIDs ...string) int {
	var svcCnt int
	p.doSync(func() {
		current, ok := p.bookings[consumerID]
		if !ok {
			return
		}
		for _, svcID := range svcIDs {
			delete(current, svcID)
		}
		if len(current) == 0 {
			delete(p.bookings, consumerID)
		} else {
			svcCnt = len(current)
			p.bookings[consumerID] = current
		}
	})
	return svcCnt
}

// Spawn runs the booked services of a consumer concurrently.
func (p *Provider) Spawn(consumerID string) {
	p.doAsync(func() {
		svcs, ok := p.bookings[consumerID]
		if !ok {
			return
		}
		svcs.Spawn()
	})
}

// doSync sends an action for execution to the backend and waits
// until its done. Right now no handling of timeouts to cover
// closed channels.
func (p *Provider) doSync(action func()) {
	donec := make(chan struct{})

	p.actionc <- func() {
		action()
		close(donec)
	}

	// Wait.
	<-donec
}

// doAsync sends an action for execution to the backend.
func (p *Provider) doAsync(action func()) {
	p.actionc <- action
}

// backend is the goroutine of the Provider.
func (p *Provider) backend() {
	defer close(p.actionc)
	for {
		select {
		case <-p.ctx.Done():
			return
		case action := <-p.actionc:
			action()
		}
	}
}
