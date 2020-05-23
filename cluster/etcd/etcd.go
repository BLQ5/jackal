/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package etcd

import (
	"time"

	v3 "github.com/coreos/etcd/clientv3"
)

func New(cfg *Config) (candidate *Leader, kv *KV, err error) {
	c, err := v3.New(v3.Config{
		DialTimeout:       time.Second * 10,
		AutoSyncInterval:  time.Second * 5,
		DialKeepAliveTime: time.Second * 5,
		Endpoints:         cfg.Endpoints,
	})
	if err != nil {
		return nil, nil, err
	}
	candidate, err = newLeader(c)
	if err != nil {
		return nil, nil, err
	}
	return candidate, &KV{cli: c}, nil
}
