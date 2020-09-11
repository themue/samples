package users

import (
	"context"
	"fmt"
)

// Registry contains all current users of the systems.
type Registry struct {
	ctx     context.Context
	actionc chan func()
	users   map[string]User
}

// StartRegistry lets the Registry run in the background.
func StartRegistry(ctx context.Context) *Registry {
	r := &Registry{
		ctx:     ctx,
		actionc: make(chan func(), 10),
		users:   make(map[string]User),
	}
	go r.backend()
	return r
}

// Create adds a new User entry.
func (r *Registry) Create(u User) error {
	var err error
	r.doSync(func() {
		if _, ok := r.users[u.ID]; ok {
			err = fmt.Errorf("user %q already exist", u.ID)
			return
		}
		r.users[u.ID] = u
	})
	return err
}

// Read retrieves a User entry by ID.
func (r *Registry) Read(id string) (User, error) {
	var u User
	var err error
	r.doSync(func() {
		ru, ok := r.users[id]
		if !ok {
			err = fmt.Errorf("user %q not found", id)
			return
		}
		u.ID = ru.ID
		u.Name = ru.Name
	})
	return u, err
}

// Update exchanges the stored User entry.
func (r *Registry) Update(u User) error {
	var err error
	r.doSync(func() {
		if _, ok := r.users[u.ID]; !ok {
			err = fmt.Errorf("user %q not found", u.ID)
			return
		}
		r.users[u.ID] = u
	})
	return err
}

// Delete removes a User entry by ID.
func (r *Registry) Delete(id string) error {
	var err error
	r.doSync(func() {
		if _, ok := r.users[id]; !ok {
			err = fmt.Errorf("user %q not found", id)
			return
		}
		delete(r.users, id)
	})
	return err
}

// doSync sends an action for execution to the backend and waits
// until its done. Right now no handling of timeouts to cover.
// closed channels.
func (r *Registry) doSync(action func()) {
	donec := make(chan struct{})

	r.actionc <- func() {
		action()
		close(donec)
	}

	<-donec
}

// backend is the goroutine of the Registry.
func (r *Registry) backend() {
	defer close(r.actionc)
	for {
		select {
		case <-r.ctx.Done():
			return
		case action := <-r.actionc:
			action()
		}
	}
}
