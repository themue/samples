// The Samples Project
//
// Copyright 2020 Frank Mueller / Oldenburg / Germany / World
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.
package consumers

import (
	"bytes"
	"context"
	"fmt"
)

// --------------------------------------------------
// Interface for Consumer Store backends.
// --------------------------------------------------

// Store defines the interface for any consumer storrage backend.
type Store interface {
	// Create adds a new Consumer.
	Create(c Consumer) error

	// Read retrieves a Consumer based on its identifier.
	Read(id string) (Consumer, error)

	// Update changes a Consumer inside the storrage.
	Update(c Consumer) error

	// Delete removes a Consumer based on its identifier.
	Delete(id string) error
}

// --------------------------------------------------
// Controller for consumer operations
// --------------------------------------------------

// Controller operates all Consumer actions like adding,
// removing, oder authenticating.
type Controller struct {
	ctx     context.Context
	actionC chan func()
	store   Store
}

// StartController starts a Consumers Controller.
func StartController(ctx context.Context, store Store) *Controller {
	cc := &Controller{
		ctx:     ctx,
		actionC: make(chan func()),
		store:   store,
	}
	go cc.backend()
	return cc
}

// Add adds a new Consumer.
func (cc *Controller) Add(c Consumer) error {
	var err error
	cc.doSync(func() {
		err = cc.store.Create(c)
	})
	if err != nil {
		return fmt.Errorf("adding consumer failed: %v", err)
	}
	return nil
}

// Read reads a Consumer by ID.
func (cc *Controller) Read(id string) (Consumer, error) {
	var c Consumer
	var err error
	cc.doSync(func() {
		c, err = cc.store.Read(id)
	})
	if err != nil {
		return Consumer{}, fmt.Errorf("reading consumer failed: %v", err)
	}
	return c, nil
}

// Remove deletes a Consumer.
func (cc *Controller) Remove(id string) {
	cc.doAsync(func() {
		cc.store.Delete(id)
	})
}

// Authenticate loads a Consumer by ID and compares the key.
func (cc *Controller) Authenticate(id string, key []byte) (Consumer, error) {
	var c Consumer
	var err error
	cc.doSync(func() {
		c, err = cc.store.Read(id)
		if err != nil {
			return
		}
		// No error, check key.
		if bytes.Compare(c.Key, key) != 0 {
			err = fmt.Errorf("key of ID %q is invalid", id)
		}
	})
	if err != nil {
		return Consumer{}, fmt.Errorf("cannot authenticate ID %q: %v", id, err)
	}
	return c, nil
}

// doSync sends an action for execution to the backend and waits
// until its done. Right now no handling of timeouts to cover
// closed channels.
func (cc *Controller) doSync(action func()) {
	doneC := make(chan struct{})

	cc.actionC <- func() {
		action()
		close(doneC)
	}

	// Wait.
	<-doneC
}

// doAsync sends an action for execution to the backend.
func (cc *Controller) doAsync(action func()) {
	cc.actionC <- action
}

// backend is the goroutine of the Controller.
func (cc *Controller) backend() {
	defer close(cc.actionC)
	for {
		select {
		case <-cc.ctx.Done():
			return
		case action := <-cc.actionC:
			action()
		}
	}
}
