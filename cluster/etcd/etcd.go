/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package etcd

import (
	"fmt"

	v3 "github.com/coreos/etcd/clientv3"
)

func New(cfg *Config) (candidate *Leader, kv *KV, err error) {
	fmt.Printf("INIT ETCD...")
	c, err := v3.New(v3.Config{Endpoints: cfg.Endpoints})
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("NEW LEADER?...")
	candidate, err = newLeader(c)
	if err != nil {
		return nil, nil, err
	}
	return candidate, &KV{cli: c}, nil
}
