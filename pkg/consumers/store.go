// The Samples Project
//
// Copyright 2020 Frank Mueller / Oldenburg / Germany / World
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.
package consumers

import (
	"context"
	"fmt"
)

// inMemoryStore contains all current consumers of the systems.
type inMemoryStore struct {
	ctx       context.Context
	actionC   chan func()
	consumers map[string]Consumer
}

// StartInMemoryStore creates a Store running simply in memory.
func StartInMemoryStore(ctx context.Context) Store {
	r := &inMemoryStore{
		ctx:       ctx,
		actionC:   make(chan func(), 1),
		consumers: make(map[string]Consumer),
	}
	go r.backend()
	return r
}

// Create adds a new Consumer entry.
func (ims *inMemoryStore) Create(c Consumer) error {
	var err error
	ims.doSync(func() {
		if _, ok := ims.consumers[c.ID]; ok {
			err = fmt.Errorf("consumer %q already exist", c.ID)
			return
		}
		ims.consumers[c.ID] = c
	})
	return err
}

// Read retrieves a Consumer entry by ID.
func (ims *inMemoryStore) Read(id string) (Consumer, error) {
	var c Consumer
	var err error
	ims.doSync(func() {
		cr, ok := ims.consumers[id]
		if !ok {
			err = fmt.Errorf("consumer %q not found", id)
			return
		}
		c = cr
	})
	return c, err
}

// Update exchanges the stored Consumer entry.
func (ims *inMemoryStore) Update(c Consumer) error {
	var err error
	ims.doSync(func() {
		if _, ok := ims.consumers[c.ID]; !ok {
			err = fmt.Errorf("consumer %q not found", c.ID)
			return
		}
		ims.consumers[c.ID] = c
	})
	return err
}

// Delete removes a Consumer entry by ID.
func (ims *inMemoryStore) Delete(id string) error {
	var err error
	ims.doSync(func() {
		if _, ok := ims.consumers[id]; !ok {
			err = fmt.Errorf("consumer %q not found", id)
			return
		}
		delete(ims.consumers, id)
	})
	return err
}

// doSync sends an action for execution to the backend and waits
// until its done. Right now no handling of timeouts or closed
// channels.
func (ims *inMemoryStore) doSync(action func()) {
	doneC := make(chan struct{})

	ims.actionC <- func() {
		action()
		close(doneC)
	}

	<-doneC
}

// backend is the goroutine of the inMemoryStore.
func (ims *inMemoryStore) backend() {
	defer close(ims.actionC)
	for {
		select {
		case <-ims.ctx.Done():
			return
		case action := <-ims.actionC:
			action()
		}
	}
}
