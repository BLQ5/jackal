/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package memorystorage

import (
	"context"
)

type Allocation struct {
	*memoryStorage
}

// NewAllocation returns an instance of Allocation in-memory storage.
func NewAllocation() *Allocation {
	return &Allocation{memoryStorage: newStorage()}
}

func (s *Allocation) RegisterAllocation(ctx context.Context, allocationID string) error {
	panic("not implemented!")
}

func (s *Allocation) UnregisterAllocation(ctx context.Context, allocationID string) error {
	panic("not implemented!")
}

func (s *Allocation) FetchAllocations(ctx context.Context) (allocationIDs []string, err error) {
	panic("not implemented!")
}
