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

	// Add the Consumers first.
	for _, add := range testData {
		err := cc.Add(add)
		if err != nil {
			t.Fatalf("adding Consumer %q failed: %v", add.ID, err)
		}
	}

	// Now read and compare.
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

	testID := "unknown"
	c, err := cc.Read(testID)
	if err == nil {
		t.Fatalf("reading unknown Consumer ID must lead to an error")
	}
	if c.ID != "" {
		t.Fatalf("read Consumer %q is not empty", testID)
	}
}

// TestRemoveConsumers verifies the removing of Consumers
// from a Controller.
func TestRemoveConsumers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	store := consumers.StartInMemoryStore(ctx)
	cc := consumers.StartController(ctx, store)

	// Add the Consumers first.
	for _, add := range testData {
		err := cc.Add(add)
		if err != nil {
			t.Fatalf("adding Consumer %q failed: %v", add.ID, err)
		}
	}

	// First read must not fail. After remove it must fail.
	_, err := cc.Read(testData[0].ID)
	if err != nil {
		t.Fatalf("reading Consumer %q failed: %v", testData[0].ID, err)
	}
	cc.Remove(testData[0].ID)
	_, err = cc.Read(testData[0].ID)
	if err == nil {
		t.Fatalf("reading Consumer %q did not failed", testData[0].ID)
	}
}

// TestAuthenticateConsumers verifies the authentication of Consumers
// by a Controller.
func TestAuthenticateConsumers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	store := consumers.StartInMemoryStore(ctx)
	cc := consumers.StartController(ctx, store)

	// Add the Consumers first.
	for _, add := range testData {
		err := cc.Add(add)
		if err != nil {
			t.Fatalf("adding Consumer %q failed: %v", add.ID, err)
		}
	}

	// Now positive and failing authentications.
	c, err := cc.Authenticate(testData[1].ID, testData[1].Key)
	if err != nil {
		t.Fatalf("authenticating Consumer %q failed: %v", testData[1], err)
	}
	if c.Name != testData[1].Name {
		t.Fatalf("authenticated Consumer %q had invalid name", testData[1])
	}
	c, err = cc.Authenticate(testData[1].ID, []byte("invalid"))
	if err == nil {
		t.Fatalf("authenticating Consumer %q did not fail", testData[1])
	}
	if c.Name != "" {
		t.Fatalf("authenticated Consumer %q is not empty", testData[1])
	}
	testID := "unknown"
	c, err = cc.Authenticate(testID, []byte("invalid"))
	if err == nil {
		t.Fatalf("authenticating Consumer %q did not fail", testID)
	}
	if c.Name != "" {
		t.Fatalf("authenticated Consumer %q is not empty", testID)
	}
}
