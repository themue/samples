// The Samples Project
//
// Copyright 2020 Frank Mueller / Oldenburg / Germany / World
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.
package consumers

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
