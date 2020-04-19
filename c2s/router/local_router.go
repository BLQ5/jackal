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

type localRouter struct {
	mu  sync.RWMutex
	tbl map[string]*resources
}

func newLocalRouter() *localRouter {
	return &localRouter{
		tbl: make(map[string]*resources),
	}
}

func (l *localRouter) route(ctx context.Context, stanza xmpp.Stanza, toJID *jid.JID) error {
	l.mu.RLock()
	rs := l.tbl[toJID.Node()]
	l.mu.RUnlock()

	if rs != nil {
		return rs.route(ctx, stanza, toJID)
	}
	return router.ErrNotAuthenticated
}

func (l *localRouter) bind(stm stream.C2S) {
	user := stm.Username()
	l.mu.RLock()
	rs := l.tbl[user]
	l.mu.RUnlock()

	if rs == nil {
		l.mu.Lock()
		rs = l.tbl[user] // avoid double initialization
		if rs == nil {
			rs = &resources{}
			l.tbl[user] = rs
		}
		l.mu.Unlock()
	}
	rs.bind(stm)
}

func (l *localRouter) unbind(user, resource string) {
	l.mu.RLock()
	rs := l.tbl[user]
	l.mu.RUnlock()

	if rs == nil {
		return
	}
	l.mu.Lock()
	rs.unbind(resource)
	if rs.len() == 0 {
		delete(l.tbl, user)
	}
	l.mu.Unlock()
}

func (l *localRouter) stream(username, resource string) stream.C2S {
	l.mu.RLock()
	rs := l.tbl[username]
	l.mu.RUnlock()

	if rs == nil {
		return nil
	}
	return rs.stream(resource)
}

func (l *localRouter) streams(username string) []stream.C2S {
	l.mu.RLock()
	rs := l.tbl[username]
	l.mu.RUnlock()

	if rs == nil {
		return nil
	}
	return rs.allStreams()
}
