/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package c2srouter

import (
	"context"
	"sync"

	"github.com/ortuman/jackal/xmpp/jid"

	"github.com/ortuman/jackal/router"
	"github.com/ortuman/jackal/stream"
	"github.com/ortuman/jackal/xmpp"
)

type resources struct {
	mu      sync.RWMutex
	streams []stream.C2S
}

func (r *resources) len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.streams)
}

func (r *resources) allStreams() []stream.C2S {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.streams
}

func (r *resources) stream(resource string) stream.C2S {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, stm := range r.streams {
		if stm.Resource() == resource {
			return stm
		}
	}
	return nil
}

func (r *resources) bind(stm stream.C2S) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := stm.Resource()
	for _, s := range r.streams {
		if s.Resource() == res {
			return
		}
	}
	r.streams = append(r.streams, stm)
}

func (r *resources) unbind(res string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, s := range r.streams {
		if s.Resource() != res {
			continue
		}
		r.streams = append(r.streams[:i], r.streams[i+1:]...)
		return
	}
}

func (r *resources) route(ctx context.Context, stanza xmpp.Stanza, toJID *jid.JID) error {
	if toJID.IsFullWithUser() {
		return r.routeToResource(ctx, stanza, toJID.Resource())
	}
	return r.broadcast(ctx, stanza)
}

func (r *resources) routeToResource(ctx context.Context, stanza xmpp.Stanza, resource string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, s := range r.streams {
		if s.Resource() != resource || !s.Presence().IsAvailable() {
			continue
		}
		s.SendElement(ctx, stanza)
		return nil
	}
	return router.ErrResourceNotFound
}

func (r *resources) broadcast(ctx context.Context, stanza xmpp.Stanza) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, s := range r.streams {
		if !s.Presence().IsAvailable() {
			continue
		}
		s.SendElement(ctx, stanza)
	}
	return nil
}
