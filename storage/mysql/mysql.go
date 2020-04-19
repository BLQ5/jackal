/*
 * Copyright (c) 2019 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // SQL driver
	"github.com/ortuman/jackal/log"
)

type Storage struct {
	*Allocation
	*User
	*Resources
	*Roster
	*Presences
	*VCard
	*Private
	*BlockList
	*PubSub
	*Offline

	h      *sql.DB
	doneCh chan chan bool
}

// New initializes MySQL storage and returns associated container.
func New(cfg *Config) (*Storage, error) {
	var err error
	host := cfg.Host
	usr := cfg.User
	pass := cfg.Password
	db := cfg.Database
	poolSize := cfg.PoolSize

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", usr, pass, host, db)
	h, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	h.SetMaxOpenConns(poolSize) // set max opened connection count

	if err := h.Ping(); err != nil {
		return nil, err
	}
	c := &Storage{
		Allocation: newAllocation(h),
		User:       newUser(h),
		Resources:  newResources(h),
		Roster:     newRoster(h),
		Presences:  newPresences(h),
		VCard:      newVCard(h),
		Private:    newPrivate(h),
		BlockList:  newBlockList(h),
		PubSub:     newPubSub(h),
		Offline:    newOffline(h),
		h:          h,
		doneCh:     make(chan chan bool, 1),
	}
	go c.loop()

	return c, nil
}

func (c *Storage) Shutdown(ctx context.Context) error {
	ch := make(chan bool)
	c.doneCh <- ch
	select {
	case <-ch:
		log.Infof("shutdown complete")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Storage) IsClusterCompatible() bool { return true }

func (c *Storage) loop() {
	tc := time.NewTicker(time.Second * 15)
	defer tc.Stop()

	for {
		select {
		case <-tc.C:
			if err := c.h.Ping(); err != nil {
				log.Error(err)
			}
		case ch := <-c.doneCh:
			if err := c.h.Close(); err != nil {
				log.Error(err)
			}
			close(ch)
			return
		}
	}
}
