// The Samples Project
//
// Copyright 2020 Frank Mueller / Oldenburg / Germany / World
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.
package consumers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/themue/samples/pkg/consumers"
)

var testData = []consumers.Consumer{
	consumers.Consumer{
		ID:   "C0",
		Key:  []byte("abcd"),
		Name: "Consumer C0",
	},
	consumers.Consumer{
		ID:   "C1",
		Key:  []byte("key"),
		Name: "Consumer C1",
	},
	consumers.Consumer{
		ID:   "C2",
		Key:  []byte("uvwxyz"),
		Name: "Consumer C2",
	},
}

// TestAddReadConsumers verifies the adding and reading of Consumers
// to a Controller.
func TestAddReadConsumers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	store := consumers.StartInMemoryStore(ctx)
	cc := consumers.StartController(ctx, store)

	for _, add := range testData {
		err := cc.Add(add)
		if err != nil {
			t.Fatalf("adding Consumer %q failed: %v", add.ID, err)
		}
	}

	for _, read := range testData {
		c, err := cc.Read(read.ID)
		if err != nil {
			t.Fatalf("reading Consumer %q failed: %v", read.ID, err)
		}
		idOK := c.ID == read.ID
		keyOK := bytes.Compare(c.Key, read.Key) == 0
		nameOK := c.Name == read.Name
		if !(idOK && keyOK && nameOK) {
			t.Fatalf("data of Consumer %q is invalid", read.ID)
		}
	}
}
