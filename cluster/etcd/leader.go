/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package etcd

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	v3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/pborman/uuid"
)

const candidateOpTimeout = time.Second * 3

type Leader struct {
	ss       *concurrency.Session
	election *concurrency.Election
	isLeader int32
	elected  int32
}

func newLeader(cli *v3.Client) (*Leader, error) {
	c := &Leader{}
	fmt.Printf("LEADER SESSION?...\n")
	ss, err := concurrency.NewSession(cli)
	if err != nil {
		return nil, err
	}
	c.ss = ss
	fmt.Printf("ELECTION...\n")
	c.election = concurrency.NewElection(c.ss, "/jackal-leader/")
	fmt.Printf("ELECTION...COMPLETED!\n")
	return c, nil
}

func (c *Leader) Elect() error {
	if !atomic.CompareAndSwapInt32(&c.elected, 0, 1) {
		return nil // already electe
	}
	go c.elect()
	return nil
}

func (c *Leader) Resign() error {
	if atomic.LoadInt32(&c.elected) == 0 {
		return nil // nothing to do here
	}
	ctx, cancel := context.WithTimeout(context.Background(), candidateOpTimeout)
	defer cancel()

	if err := c.election.Resign(ctx); err != nil {
		return err
	}
	atomic.StoreInt32(&c.isLeader, 0)

	return c.ss.Close()
}

func (c *Leader) IsLeader() bool {
	return atomic.LoadInt32(&c.isLeader) == 1
}

func (c *Leader) elect() {
	val := uuid.New()
	if err := c.election.Campaign(context.Background(), val); err != nil {
		log.Fatalf("%v", err)
		return
	}
	atomic.StoreInt32(&c.isLeader, 1)
}
