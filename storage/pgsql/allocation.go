/*
 * Copyright (c) 2018 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package pgsql

import (
	"context"
	"database/sql"
)

type Allocation struct {
	*pgSQLStorage
}

func newAllocation(db *sql.DB) *Allocation {
	return &Allocation{
		pgSQLStorage: newStorage(db),
	}
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
