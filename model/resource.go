/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package model

import (
	"bytes"
	"encoding/gob"

	"github.com/ortuman/jackal/xmpp/jid"
)

type Resource struct {
	AllocationID string
	JID          *jid.JID
	Priority     int8
}

// FromBytes deserializes a Resource entity from its binary representation.
func (r *Resource) FromBytes(buf *bytes.Buffer) error {
	var jStr string

	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&r.AllocationID); err != nil {
		return err
	}
	if err := dec.Decode(&jStr); err != nil {
		return err
	}
	r.JID, _ = jid.NewWithString(jStr, true)

	return dec.Decode(&r.Priority)
}

// ToBytes converts a Resource entity to its binary representation.
func (r *Resource) ToBytes(buf *bytes.Buffer) error {
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(r.AllocationID); err != nil {
		return err
	}
	if err := enc.Encode(r.JID.String()); err != nil {
		return err
	}
	return enc.Encode(r.Priority)
}
