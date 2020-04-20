/*
 * Copyright (c) 2018 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package c2s

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/ortuman/jackal/cluster"
	"github.com/ortuman/jackal/component"
	"github.com/ortuman/jackal/log"
	"github.com/ortuman/jackal/module"
	"github.com/ortuman/jackal/router"
	"github.com/ortuman/jackal/storage"
	"github.com/pkg/errors"
)

const houseKeepingInterval = time.Second * 3

const (
	streamNamespace           = "http://etherx.jabber.org/streams"
	tlsNamespace              = "urn:ietf:params:xml:ns:xmpp-tls"
	compressProtocolNamespace = "http://jabber.org/protocol/compress"
	bindNamespace             = "urn:ietf:params:xml:ns:xmpp-bind"
	sessionNamespace          = "urn:ietf:params:xml:ns:xmpp-session"
	saslNamespace             = "urn:ietf:params:xml:ns:xmpp-sasl"
	blockedErrorNamespace     = "urn:xmpp:blocking:errors"
)

type c2sServer interface {
	start()
	shutdown(ctx context.Context) error
}

var createC2SServer = newC2SServer

// C2S represents a client-to-server connection manager.
type C2S struct {
	servers      map[string]c2sServer
	allocationSt storage.Allocation
	cluster      *cluster.Cluster
	closeCh      chan chan struct{}
	started      uint32
}

// New returns a new instance of a c2s connection manager.
func New(
	configs []Config,
	mods *module.Modules,
	comps *component.Components,
	router router.Router,
	userSt storage.User,
	resourcesSt storage.Resources,
	blockListSt storage.BlockList,
	allocationSt storage.Allocation,
	cluster *cluster.Cluster,
) (*C2S, error) {
	if len(configs) == 0 {
		return nil, errors.New("at least one c2s configuration is required")
	}
	c := &C2S{
		servers:      make(map[string]c2sServer),
		closeCh:      make(chan chan struct{}, 1),
		allocationSt: allocationSt,
		cluster:      cluster,
	}
	for _, config := range configs {
		srv := createC2SServer(&config, mods, comps, router, userSt, resourcesSt, blockListSt)
		c.servers[config.ID] = srv
	}
	return c, nil
}

// Start initializes c2s manager spawning every single server.
func (c *C2S) Start() error {
	if atomic.CompareAndSwapUint32(&c.started, 0, 1) {
		for _, srv := range c.servers {
			go srv.start()
		}
		if c.cluster != nil {
			// do master election
			if err := c.cluster.Elect(); err != nil {
				return err
			}
			// join to cluster
			if err := c.cluster.Join(); err != nil {
				return err
			}
			go c.loop()
		}
		log.Infof("c2s started")
	}
	return nil
}

// Shutdown gracefully shuts down c2s manager.
func (c *C2S) Shutdown(ctx context.Context) error {
	if atomic.CompareAndSwapUint32(&c.started, 1, 0) {
		for _, srv := range c.servers {
			if err := srv.shutdown(ctx); err != nil {
				return err
			}
		}
		ch := make(chan struct{})
		c.closeCh <- ch
		select {
		case <-ch:
			break
		case <-ctx.Done():
			return ctx.Err()
		}
		log.Infof("shutdown complete")
	}
	return nil
}

func (c *C2S) loop() {
	tc := time.NewTicker(houseKeepingInterval)
	defer tc.Stop()

	for {
		select {
		case <-tc.C:
			if err := c.houseKeeping(); err != nil {
				log.Warnf("housekeeping task error: %v", err)
			}
		case ch := <-c.closeCh:
			close(ch)
			return
		}
	}
}

func (c *C2S) houseKeeping() error {
	if c.cluster != nil && !c.cluster.IsLeader() {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), (houseKeepingInterval*5)/10)
	defer cancel()

	allocations, err := c.allocationSt.FetchAllocations(ctx)
	if err != nil {
		return err
	}
	members := c.cluster.Members()
	for _, alloc := range allocations {
		if m := members.Member(alloc.ID); m != nil {
			continue
		}
		// unregister inactive allocations
		if err := c.allocationSt.UnregisterAllocation(ctx, alloc.ID); err != nil {
			return err
		}
		log.Infof("unregistered dangling allocation: %s", alloc.ID)
	}
	return nil
}
