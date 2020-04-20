/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package c2srouter

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/ortuman/jackal/cluster"
	"github.com/ortuman/jackal/log"
	"github.com/ortuman/jackal/util/pool"
	"github.com/ortuman/jackal/xmpp"
	"github.com/ortuman/jackal/xmpp/jid"
	"github.com/sony/gobreaker"
	"golang.org/x/net/http2"
)

type clusterRouter struct {
	hClient    *http.Client
	cb         *gobreaker.CircuitBreaker
	pool       *pool.BufferPool
	memberList cluster.MemberList
}

func newClusterRouter(memberList cluster.MemberList) (*clusterRouter, error) {
	h2cTransport := &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}
	r := &clusterRouter{
		hClient:    &http.Client{Transport: h2cTransport},
		cb:         gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		pool:       pool.NewBufferPool(),
		memberList: memberList,
	}
	return r, nil
}

func (r *clusterRouter) route(ctx context.Context, stanza xmpp.Stanza, toJIDs []*jid.JID, allocationID string) error {
	member := r.memberList.Members().Member(allocationID)
	if member == nil {
		log.Warnf(fmt.Sprintf("allocation %s not found", allocationID))
		return nil
	}
	buf := r.pool.Get()
	defer r.pool.Put(buf)

	if err := stanza.ToXML(buf, true); err != nil {
		return err
	}
	var toParam strings.Builder
	for i, toJID := range toJIDs {
		if i != 0 {
			toParam.WriteString(",")
		}
		toParam.WriteString(toJID.String())
	}
	reqURL := fmt.Sprintf("http://%s:%s/route?to=%s", member.Host, member.Port, url.QueryEscape(toParam.String()))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/xml")

	_, err = r.cb.Execute(func() (i interface{}, e error) {
		resp, err := r.hClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("response status code: %d", resp.StatusCode)
		}
		return nil, nil
	})
	return err
}
