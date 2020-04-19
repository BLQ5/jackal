/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package storage

import "context"

type Allocation interface {
	RegisterAllocation(ctx context.Context, allocationID string) error

	UnregisterAllocation(ctx context.Context, allocationID string) error

	FetchAllocations(ctx context.Context) (allocationIDs []string, err error)
}
