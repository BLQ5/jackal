/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package memorystorage

import (
	"context"
	"strings"

	"github.com/ortuman/jackal/model/serializer"

	"github.com/ortuman/jackal/model"
)

type Allocation struct {
	*memoryStorage
}

// NewAllocation returns an instance of Allocation in-memory storage.
func NewAllocation() *Allocation {
	return &Allocation{memoryStorage: newStorage()}
}

func (m *Allocation) RegisterAllocation(_ context.Context, allocation *model.Allocation) error {
	return m.saveEntity(allocationKey(allocation.ID), allocation)
}

func (m *Allocation) UnregisterAllocation(_ context.Context, allocationID string) error {
	return m.deleteKey(allocationKey(allocationID))
}

func (m *Allocation) FetchAllocations(_ context.Context) ([]model.Allocation, error) {
	var allocations []model.Allocation
	if err := m.inReadLock(func() error {
		for k, b := range m.b {
			if !strings.HasPrefix(k, "allocations:") {
				continue
			}
			var alloc model.Allocation
			if err := serializer.Deserialize(b, &alloc); err != nil {
				return err
			}
			allocations = append(allocations, alloc)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return allocations, nil
}

func allocationKey(allocationID string) string {
	return "allocations:" + allocationID
}
