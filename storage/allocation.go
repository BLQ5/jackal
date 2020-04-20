/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package storage

import (
	"context"

	"github.com/ortuman/jackal/model"
)

type Allocation interface {
	RegisterAllocation(ctx context.Context, allocation *model.Allocation) error

	UnregisterAllocation(ctx context.Context, allocationID string) error

	FetchAllocations(ctx context.Context) (allocations []model.Allocation, err error)
}
