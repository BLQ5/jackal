/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package model

import (
	"bytes"
	"encoding/gob"
)

// Allocations represents a cluster allocation instance.
type Allocation struct {
	// ID represents an allocation instance unique identifier
	ID string
}

// FromBytes deserializes a Allocation entity from its binary representation.
func (a *Allocation) FromBytes(buf *bytes.Buffer) error {
	dec := gob.NewDecoder(buf)
	return dec.Decode(&a.ID)
}

// ToBytes converts a Allocation entity to its binary representation.
func (a *Allocation) ToBytes(buf *bytes.Buffer) error {
	enc := gob.NewEncoder(buf)
	return enc.Encode(a.ID)
}
