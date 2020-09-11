package metaweather_test

import (
	"context"
	"testing"
	"time"

	"github.com/themue/training-samples/pkg/metaweather"
)

func TestSubscribe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sub := metaweather.StartSubscriber(ctx, 10*time.Second)

	cities := sub.Subscribe("london")
	if len(cities) != 1 {
		t.Fatalf("illegal number of cities: %v", cities)
	}
	cities = sub.Subscribe("san")
	if len(cities) != 11 {
		t.Fatalf("illegal number of cities: %v", cities)
	}
	cities = sub.Subscribe("thisissomestrangelocationwhichdoesnotexist")
	if len(cities) != 0 {
		t.Fatalf("illegal number of cities: %v", cities)
	}
}

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
