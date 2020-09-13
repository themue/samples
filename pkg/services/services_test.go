package services_test

import (
	"context"
	"sync"
	"testing"

	"github.com/themue/training-samples/pkg/services"
)

// TestSpawnServices validates the correct execution
// of a number of services.
func TestSpawnServices(t *testing.T) {
	var wg sync.WaitGroup
	var dsia, dsib, dsic int
	cba := func(i int) {
		dsia = i
	}
	cbb := func(i int) {
		dsib = i
	}
	cbc := func(i int) {
		dsic = i
	}
	svcs := services.Services{
		"a": newDummyService("a", cba, &wg),
		"b": newDummyService("b", cbb, &wg),
		"c": newDummyService("c", cbc, &wg),
	}

	wg.Add(9)
	svcs.Spawn()
	svcs.Spawn()
	svcs.Spawn()
	wg.Wait()

	if dsia != 3 {
		t.Fatalf("a has wrong value: %d", dsia)
	}
	if dsib != 3 {
		t.Fatalf("b has wrong value: %d", dsib)
	}
	if dsic != 3 {
		t.Fatalf("c has wrong value: %d", dsib)
	}
}

// TestProviderBookUnbook validates booking and unbooking
// of services.
func TestProviderBookUnbook(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	p := services.StartProvider(ctx)
	svca := newDummyService("a", func(i int) {}, &wg)
	svcb := newDummyService("b", func(i int) {}, &wg)
	svcc := newDummyService("c", func(i int) {}, &wg)

	svcCnt := p.Book("foo", svca, svcb)
	if svcCnt != 2 {
		t.Fatalf("invalid number of services, expect 2: %d", svcCnt)
	}
	svcCnt = p.Book("foo", svcc)
	if svcCnt != 3 {
		t.Fatalf("invalid number of services, expect 3: %d", svcCnt)
	}
	svcCnt = p.Book("foo", svca)
	if svcCnt != 3 {
		t.Fatalf("invalid number of services, expect 3: %d", svcCnt)
	}
	svcCnt = p.Unbook("foo", "a", "b")
	if svcCnt != 1 {
		t.Fatalf("invalid number of services, expect 1: %d", svcCnt)
	}
	svcCnt = p.Unbook("foo", "a", "c")
	if svcCnt != 0 {
		t.Fatalf("invalid number of services, expect 0: %d", svcCnt)
	}
	svcCnt = p.Unbook("bar", "b")
	if svcCnt != 0 {
		t.Fatalf("invalid number of services, expect 0: %d", svcCnt)
	}
}

// TestProviderSpawn validates spawning booked services via
// the ServiceProvider.
func TestProviderSpawn(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p := services.StartProvider(ctx)
	var wg sync.WaitGroup
	var dsia, dsib, dsic int
	cba := func(i int) {
		dsia = i
	}
	cbb := func(i int) {
		dsib = i
	}
	cbc := func(i int) {
		dsic = i
	}
	svca := newDummyService("a", cba, &wg)
	svcb := newDummyService("b", cbb, &wg)
	svcc := newDummyService("c", cbc, &wg)

	svcCnt := p.Book("foo", svca, svcb, svcc)
	if svcCnt != 3 {
		t.Fatalf("invalid number of services, expect 3: %d", svcCnt)
	}

	wg.Add(9)
	p.Spawn("foo")
	p.Spawn("foo")
	p.Spawn("foo")
	p.Spawn("bar")
	wg.Wait()

	if dsia != 3 {
		t.Fatalf("a has wrong value: %d", dsia)
	}
	if dsib != 3 {
		t.Fatalf("b has wrong value: %d", dsib)
	}
	if dsic != 3 {
		t.Fatalf("c has wrong value: %d", dsib)
	}
}

// -----
// dummyService is a simple implementation
// of Service for testing purposes.
// -----

type dummyService struct {
	mu       sync.Mutex
	id       string
	incr     int
	callback func(int)
	wg       *sync.WaitGroup
}

func newDummyService(id string, cb func(int), wg *sync.WaitGroup) services.Service {
	return &dummyService{
		id:       id,
		callback: cb,
		wg:       wg,
	}
}

func (s *dummyService) ID() string {
	return s.id
}

func (s *dummyService) Do() error {
	s.incr++
	s.callback(s.incr)
	s.wg.Done()
	return nil
}
