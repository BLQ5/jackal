/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package model

import (
	"bytes"
	"encoding/gob"

	capsmodel "github.com/ortuman/jackal/model/capabilities"
	"github.com/ortuman/jackal/xmpp"
)

// ExtPresence represents an extended presence type.
type ExtPresence struct {
	// Presence contains presence stanza value.
	Presence *xmpp.Presence

	// Caps contains presence associates capabilities.
	Caps *capsmodel.Capabilities
}

// FromBytes deserializes a ExtPresence entity from its binary representation.
func (p *ExtPresence) FromBytes(buf *bytes.Buffer) error {
	presence, err := xmpp.NewPresenceFromBytes(buf)
	if err != nil {
		return err
	}
	var hasCaps bool

	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&hasCaps); err != nil {
		return err
	}
	p.Presence = presence
	if hasCaps {
		return dec.Decode(&p.Caps)
	}
	return nil
}

// ToBytes converts a ExtPresence entity to its binary representation.
func (p *ExtPresence) ToBytes(buf *bytes.Buffer) error {
	if err := p.Presence.ToBytes(buf); err != nil {
		return err
	}
	enc := gob.NewEncoder(buf)

	hasCaps := p.Caps != nil
	if err := enc.Encode(hasCaps); err != nil {
		return err
	}
	if p.Caps != nil {
		if err := enc.Encode(p.Caps); err != nil {
			return err
		}
	}
	return nil
}
