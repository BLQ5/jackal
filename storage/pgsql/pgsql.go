/*
 * Copyright (c) 2019 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/ortuman/jackal/log"
)

// pingInterval defines how often to check the connection
var pingInterval = 15 * time.Second

// pingTimeout defines how long to wait for pong from server
var pingTimeout = 10 * time.Second

type Storage struct {
	*Allocation
	*User
	*Resources
	*Presences
	*Roster
	*VCard
	*Private
	*BlockList
	*PubSub
	*Offline

	h          *sql.DB
	cancelPing context.CancelFunc
	doneCh     chan chan bool
}

// New initializes PgSQL storage and returns associated container.
func New(cfg *Config) (*Storage, error) {
	var err error

	sq.StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Database, cfg.SSLMode)

	h, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	h.SetMaxOpenConns(cfg.PoolSize) // set max opened connection count

	ctx, cancel := context.WithCancel(context.Background())

	c := &Storage{
		Allocation: newAllocation(h),
		User:       newUser(h),
		Resources:  newResources(h),
		Presences:  newPresences(h),
		Roster:     newRoster(h),
		VCard:      newVCard(h),
		Private:    newPrivate(h),
		BlockList:  newBlockList(h),
		PubSub:     newPubSub(h),
		Offline:    newOffline(h),
		h:          h,
		cancelPing: cancel,
		doneCh:     make(chan chan bool, 1),
	}
	if err := c.ping(context.Background()); err != nil {
		return nil, err
	}
	go c.loop(ctx)

	return c, nil
}

func (c *Storage) Shutdown(ctx context.Context) error {
	ch := make(chan bool)
	c.doneCh <- ch
	select {
	case <-ch:
		log.Infof("Shutdown complete")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Storage) IsClusterCompatible() bool { return true }

func (c *Storage) loop(ctx context.Context) {
	tick := time.NewTicker(pingInterval)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			if err := c.ping(ctx); err != nil {
				log.Error(err)
			}

		case ch := <-c.doneCh:
			if err := c.h.Close(); err != nil {
				log.Error(err)
			}
			close(ch)
			return

		case <-ctx.Done():
			return
		}
	}
}

// ping sends a ping request to the server and outputs any error to log
func (c *Storage) ping(ctx context.Context) error {
	pingCtx, cancel := context.WithDeadline(ctx, time.Now().Add(pingTimeout))
	defer cancel()

	return c.h.PingContext(pingCtx)
}
