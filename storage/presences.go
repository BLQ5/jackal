/*
 * Copyright (c) 2018 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package storage

import (
	"context"

	"github.com/ortuman/jackal/model"
	capsmodel "github.com/ortuman/jackal/model/capabilities"
	"github.com/ortuman/jackal/xmpp"
	"github.com/ortuman/jackal/xmpp/jid"
)

type Presences interface {
	// UpsertPresence inserts or updates a presence and links it to certain allocation.
	// On insertion 'inserted' return parameter will be true.
	UpsertPresence(ctx context.Context, presence *xmpp.Presence, jid *jid.JID, allocationID string) (inserted bool, err error)

	// FetchPresence retrieves from storage a previously registered presence.
	FetchPresence(ctx context.Context, jid *jid.JID) (*model.ExtPresence, error)

	// FetchPresencesMatchingJID retrives all storage presences matching a certain JID
	FetchPresencesMatchingJID(ctx context.Context, jid *jid.JID) ([]model.ExtPresence, error)

	// DeletePresence removes from storage a concrete registered presence.
	DeletePresence(ctx context.Context, jid *jid.JID) error

	// UpsertCapabilities inserts capabilities associated to a node+ver pair, or updates them if previously inserted..
	UpsertCapabilities(ctx context.Context, caps *capsmodel.Capabilities) error

	// FetchCapabilities fetches capabilities associated to a give node and ver.
	FetchCapabilities(ctx context.Context, node, ver string) (*capsmodel.Capabilities, error)
}
