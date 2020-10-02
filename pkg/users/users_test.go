// The Samples Project
//
// Copyright 2020 Frank Mueller / Oldenburg / Germany / World
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.
package users_test

import (
	"context"
	"testing"

	"github.com/themue/training-samples/pkg/users"
)

func TestCreateRead(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ureg := users.StartRegistry(ctx)

	uFooIn := users.User{"foo", "A. Foo"}
	uBarIn := users.User{"bar", "B. Bar"}

	if err := ureg.Create(uFooIn); err != nil {
		t.Fatalf("creating %v failed: %v", uFooIn, err)
	}
	if err := ureg.Create(uFooIn); err == nil {
		t.Fatalf("%v had been created twice", uFooIn)
	}
	if err := ureg.Create(uBarIn); err != nil {
		t.Fatalf("creating %v failed: %v", uBarIn, err)
	}

	uFooOut, err := ureg.Read("foo")
	if err != nil {
		t.Fatalf("reading %q failed: %v", "foo", err)
	}
	if uFooOut.ID != uFooIn.ID {
		t.Fatalf("user %v had wrong ID", uFooOut)
	}
	uBarOut, err := ureg.Read("bar")
	if err != nil {
		t.Fatalf("reading %q failed: %v", "bar", err)
	}
	if uBarOut.ID != uBarIn.ID {
		t.Fatalf("user %v had wrong ID", uBarOut)
	}
	uBazOut, err := ureg.Read("baz")
	if err == nil {
		t.Fatalf("read invalid %v", uBazOut)
	}
}

func TestUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ureg := users.StartRegistry(ctx)

	ureg.Create(users.User{"foo", "A. Foo"})
	ureg.Create(users.User{"bar", "B. Bar"})

	ureg.Update(users.User{"foo", "A. Bar"})
	ureg.Update(users.User{"bar", "B. Foo"})

	uFooOut, err := ureg.Read("foo")
	if err != nil {
		t.Fatalf("reading %q failed: %v", "foo", err)
	}
	if uFooOut.Name != "A. Bar" {
		t.Fatalf("user %v had wrong name", uFooOut)
	}
	uBarOut, err := ureg.Read("bar")
	if err != nil {
		t.Fatalf("reading %q failed: %v", "bar", err)
	}
	if uBarOut.Name != "B. Foo" {
		t.Fatalf("user %v had wrong name", uBarOut)
	}
}

func TestDelete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ureg := users.StartRegistry(ctx)

	ureg.Create(users.User{"foo", "A. Foo"})
	_, err := ureg.Read("foo")
	if err != nil {
		t.Fatalf("reading %q failed: %v", "foo", err)
	}
	err = ureg.Delete("foo")
	if err != nil {
		t.Fatalf("deleting %q failed: %v", "foo", err)
	}
	_, err = ureg.Read("foo")
	if err == nil {
		t.Fatalf("r%q had not been deleted", "foo")
	}
}
