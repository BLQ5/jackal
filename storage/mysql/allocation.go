/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package mysql

import (
	"context"
	"database/sql"
)

type Allocation struct {
	*mySQLStorage
}

func newAllocation(db *sql.DB) *Allocation {
	return &Allocation{
		mySQLStorage: newStorage(db),
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
