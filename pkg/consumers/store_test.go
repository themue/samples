// The Samples Project
//
// Copyright 2020 Frank Mueller / Oldenburg / Germany / World
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.
package consumers_test

import (
	"context"
	"testing"

	"github.com/themue/samples/pkg/consumers"
)

// TestCreateRead verifies create and read operations inside an in-memory store.
func TestCreateRead(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	store := consumers.StartInMemoryStore(ctx)

	cFooIn := consumers.Consumer{"foo", "A. Foo"}
	cBarIn := consumers.Consumer{"bar", "B. Bar"}

	if err := store.Create(cFooIn); err != nil {
		t.Fatalf("creating %v failed: %v", cFooIn, err)
	}
	if err := store.Create(cFooIn); err == nil {
		t.Fatalf("%v had been created twice", cFooIn)
	}
	if err := store.Create(cBarIn); err != nil {
		t.Fatalf("creating %v failed: %v", cBarIn, err)
	}

	cFooOut, err := store.Read("foo")
	if err != nil {
		t.Fatalf("reading %q failed: %v", "foo", err)
	}
	if cFooOut.ID != cFooIn.ID {
		t.Fatalf("user %v had wrong ID", cFooOut)
	}
	cBarOut, err := store.Read("bar")
	if err != nil {
		t.Fatalf("reading %q failed: %v", "bar", err)
	}
	if cBarOut.ID != cBarIn.ID {
		t.Fatalf("user %v had wrong ID", cBarOut)
	}
	cBazOut, err := store.Read("baz")
	if err == nil {
		t.Fatalf("read invalid %v", cBazOut)
	}
}

// TestUpdate verifies the updating of a Consumer inside an in-memory store.
func TestUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	store := consumers.StartInMemoryStore(ctx)

	store.Create(consumers.Consumer{"foo", "A. Foo"})
	store.Create(consumers.Consumer{"bar", "B. Bar"})

	store.Update(consumers.Consumer{"foo", "A. Bar"})
	store.Update(consumers.Consumer{"bar", "B. Foo"})

	cFooOut, err := store.Read("foo")
	if err != nil {
		t.Fatalf("reading %q failed: %v", "foo", err)
	}
	if cFooOut.Name != "A. Bar" {
		t.Fatalf("consumer %v had wrong name", cFooOut)
	}
	cBarOut, err := store.Read("bar")
	if err != nil {
		t.Fatalf("reading %q failed: %v", "bar", err)
	}
	if cBarOut.Name != "B. Foo" {
		t.Fatalf("consumer %v had wrong name", cBarOut)
	}
}

// TestDelete verifies the removing of a Consumer from an in-memory store.
func TestDelete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	store := consumers.StartInMemoryStore(ctx)

	store.Create(consumers.Consumer{"foo", "A. Foo"})
	_, err := store.Read("foo")
	if err != nil {
		t.Fatalf("reading %q failed: %v", "foo", err)
	}
	err = store.Delete("foo")
	if err != nil {
		t.Fatalf("deleting %q failed: %v", "foo", err)
	}
	_, err = store.Read("foo")
	if err == nil {
		t.Fatalf("consumer %q had not been deleted", "foo")
	}
}
