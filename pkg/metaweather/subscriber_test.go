package metaweather_test

import (
	"context"
	"testing"
	"time"

	"github.com/themue/training-samples/pkg/metaweather"
)

// TestSubscribe verifies the subscription to a number of
// locations together with the correct returned locations.
func TestSubscribe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sub := metaweather.StartSubscriber(ctx, 10*time.Second)

	locations := sub.Subscribe("london")
	if len(locations) != 1 {
		t.Fatalf("illegal number of locations: %v", locations)
	}
	locations = sub.Subscribe("san")
	if len(locations) != 11 {
		t.Fatalf("illegal number of locations: %v", locations)
	}
	locations = sub.Subscribe("thisissomestrangelocationwhichdoesnotexist")
	if len(locations) != 0 {
		t.Fatalf("illegal number of locations: %v", locations)
	}
}

// TestUpdates verifies the background update of subscribed locations.
func TestUpdates(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sub := metaweather.StartSubscriber(ctx, 5*time.Second)

	sub.Subscribe("london")

	weathers := sub.Fetch("london")
	if len(weathers) != 1 {
		t.Fatalf("illegal number of cities: %v", weathers)
	}
	t1 := weathers[0].Time
	time.Sleep(20 * time.Second)
	weathers = sub.Fetch("london")
	t2 := weathers[0].Time
	if t1 == t2 {
		t.Fatalf("time #1 (%v) has to be different from time #2 (%v)", t1, t2)
	}
}
