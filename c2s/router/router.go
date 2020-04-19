/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package c2srouter

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/ortuman/jackal/cluster"
	"github.com/ortuman/jackal/log"
	"github.com/ortuman/jackal/model"
	"github.com/ortuman/jackal/router"
	"github.com/ortuman/jackal/storage"
	"github.com/ortuman/jackal/stream"
	"github.com/ortuman/jackal/xmpp"
	"github.com/ortuman/jackal/xmpp/jid"
)

type c2sRouter struct {
	userSt        storage.User
	blockListSt   storage.BlockList
	resourcesSt   storage.Resources
	cluster       *cluster.Cluster
	localRouter   *localRouter
	clusterRouter *clusterRouter
	allocationSt  storage.Allocation
	closeCh       chan chan struct{}
}

func New(
	userSt storage.User,
	resourcesSt storage.Resources,
	blockListSt storage.BlockList,
	allocationSt storage.Allocation,
	cluster *cluster.Cluster,
) (router.C2SRouter, error) {
	r := &c2sRouter{
		userSt:       userSt,
		blockListSt:  blockListSt,
		resourcesSt:  resourcesSt,
		localRouter:  newLocalRouter(),
		cluster:      cluster,
		allocationSt: allocationSt,
		closeCh:      make(chan chan struct{}, 1),
	}
	if cluster != nil {
		clusterRouter, err := newClusterRouter(cluster.MemberList)
		if err != nil {
			return nil, err
		}
		r.clusterRouter = clusterRouter

		// register local router as cluster stanza handler
		cluster.RegisterStanzaHandler(r.localRouter.route)

		if err := r.cluster.Elect(); err != nil {
			return nil, err
		}
		if err := r.cluster.Join(); err != nil {
			return nil, err
		}
		go r.loop()
	}
	return r, nil
}

func (r *c2sRouter) Route(ctx context.Context, stanza xmpp.Stanza, validations router.C2SRoutingValidations) error {
	fromJID := stanza.FromJID()
	toJID := stanza.ToJID()

	// apply validations
	username := stanza.ToJID().Node()
	if (validations & router.UserExistence) > 0 {
		exists, err := r.userSt.UserExists(ctx, username) // user exists?
		if err != nil {
			return err
		}
		if !exists {
			return router.ErrNotExistingAccount
		}
	}
	if (validations & router.BlockedDestinationJID) > 0 {
		if r.isBlockedJID(ctx, toJID, fromJID.Node()) { // check whether destination JID is blocked
			return router.ErrBlockedJID
		}
	}
	// fetch available resources
	resources, err := r.resourcesSt.FetchResources(ctx, toJID.Node(), toJID.Domain())
	if err != nil {
		return err
	}
	if len(resources) == 0 {
		return router.ErrNotAuthenticated
	}
	return r.route(ctx, stanza, resources)
}

func (r *c2sRouter) Bind(stm stream.C2S) {
	r.localRouter.bind(stm)

	log.Infof("bound c2s stream... (%s/%s)", stm.Username(), stm.Resource())
}

func (r *c2sRouter) Unbind(user, resource string) {
	r.localRouter.unbind(user, resource)

	log.Infof("unbound c2s stream... (%s/%s)", user, resource)
}

func (r *c2sRouter) Stream(username, resource string) stream.C2S {
	return r.localRouter.stream(username, resource)
}

func (r *c2sRouter) Streams(username string) []stream.C2S {
	return r.localRouter.streams(username)
}

func (r *c2sRouter) route(ctx context.Context, stanza xmpp.Stanza, resources []model.Resource) error {
	toJID := stanza.ToJID()
	if toJID.IsFullWithUser() {
		return r.routeToFullResource(ctx, stanza, resources)
	}
	switch msg := stanza.(type) {
	case *xmpp.Message:
		routed, err := r.routeToPrioritaryResources(ctx, msg, resources)
		if err != nil {
			return err
		}
		if !routed {
			goto route2all
		}
		return nil
	}
route2all:
	return r.routeToAllResources(ctx, stanza, resources)
}

func (r *c2sRouter) routeToFullResource(ctx context.Context, stanza xmpp.Stanza, resources []model.Resource) error {
	toJID := stanza.ToJID()
	for _, res := range resources {
		if stanza.ToJID().Resource() != res.JID.Resource() {
			continue
		}
		return r.routeToAllocation(ctx, stanza, []*jid.JID{toJID}, res.AllocationID)
	}
	return router.ErrResourceNotFound
}

func (r *c2sRouter) routeToPrioritaryResources(ctx context.Context, stanza xmpp.Stanza, resources []model.Resource) (routed bool, err error) {
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Priority > resources[j].Priority
	})
	highestPriority := resources[0].Priority
	if highestPriority <= 0 {
		return false, nil // no prioritary presence found
	}
	var prioritaryResources []model.Resource
	for _, res := range resources {
		if res.Priority != highestPriority {
			break
		}
		prioritaryResources = append(prioritaryResources, res)
	}
	// broacast to prioritary resources
	if err := r.routeToAllResources(ctx, stanza, prioritaryResources); err != nil {
		return false, err
	}
	return true, nil
}

func (r *c2sRouter) routeToAllResources(ctx context.Context, stanza xmpp.Stanza, resources []model.Resource) error {
	routeTbl := make(map[string][]*jid.JID)
	for _, res := range resources {
		routeTbl[res.AllocationID] = append(routeTbl[res.AllocationID], res.JID)
	}
	errCh := make(chan error, len(routeTbl))

	var wg sync.WaitGroup
	for k, v := range routeTbl {
		wg.Add(1)

		go func(allocationID string, toJIDs []*jid.JID) {
			defer wg.Done()
			if err := r.routeToAllocation(ctx, stanza, toJIDs, allocationID); err != nil {
				errCh <- err
			}
		}(k, v)
	}
	go func() {
		wg.Wait()
		errCh <- nil
	}()
	return <-errCh
}

func (r *c2sRouter) routeToAllocation(ctx context.Context, stanza xmpp.Stanza, toJIDs []*jid.JID, allocID string) error {
	if r.clusterRouter == nil || r.cluster.IsLocalAllocationID(allocID) {
		for _, toJID := range toJIDs {
			return r.localRouter.route(ctx, stanza, toJID)
		}
	}
	return r.clusterRouter.route(ctx, stanza, toJIDs, allocID)
}

func (r *c2sRouter) isBlockedJID(ctx context.Context, j *jid.JID, username string) bool {
	blockList, err := r.blockListSt.FetchBlockListItems(ctx, username)
	if err != nil {
		log.Error(err)
		return false
	}
	if len(blockList) == 0 {
		return false
	}
	blockListJIDs := make([]jid.JID, len(blockList))
	for i, listItem := range blockList {
		j, _ := jid.NewWithString(listItem.JID, true)
		blockListJIDs[i] = *j
	}
	for _, blockedJID := range blockListJIDs {
		if blockedJID.Matches(j) {
			return true
		}
	}
	return false
}

func (r *c2sRouter) shutdown(ctx context.Context) error {
	ch := make(chan struct{})
	r.closeCh <- ch
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *c2sRouter) loop() {
	tc := time.NewTicker(houseKeepingInterval)
	defer tc.Stop()

	for {
		select {
		case <-tc.C:
			if err := r.houseKeeping(); err != nil {
				log.Warnf("housekeeping task error: %v", err)
			}
		case ch := <-r.closeCh:
			close(ch)
			return
		}
	}
}

func (r *c2sRouter) houseKeeping() error {
	if !r.cluster.IsLeader() {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), (houseKeepingInterval*5)/10)
	defer cancel()

	allocIDs, err := r.allocationSt.FetchAllocations(ctx)
	if err != nil {
		return err
	}
	members := r.cluster.Members()
	for _, allocID := range allocIDs {
		if m := members.Member(allocID); m != nil {
			continue
		}
		// unregister inactive allocations
		if err := r.allocationSt.UnregisterAllocation(ctx, allocID); err != nil {
			return err
		}
	}
	return nil
}
