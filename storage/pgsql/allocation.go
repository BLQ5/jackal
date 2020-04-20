/*
 * Copyright (c) 2018 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package pgsql

import (
	"context"
	"database/sql"

	"github.com/ortuman/jackal/model"

	sq "github.com/Masterminds/squirrel"
)

type Allocation struct {
	*pgSQLStorage
}

func newAllocation(db *sql.DB) *Allocation {
	return &Allocation{
		pgSQLStorage: newStorage(db),
	}
}

func (s *Allocation) RegisterAllocation(ctx context.Context, allocation *model.Allocation) error {
	_, err := sq.Insert("allocations").
		Columns("allocation_id").
		Suffix("ON CONFLICT (allocation_id) DO NOTHING").
		Values(allocation.ID).
		RunWith(s.db).ExecContext(ctx)
	return err
}

func (s *Allocation) UnregisterAllocation(ctx context.Context, allocationID string) error {
	return s.inTransaction(ctx, func(tx *sql.Tx) error {
		_, err := sq.Delete("presences").Where(sq.Eq{"allocation_id": allocationID}).
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return err
		}
		_, err = sq.Delete("resources").Where(sq.Eq{"allocation_id": allocationID}).
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return err
		}
		_, err = sq.Delete("allocations").Where(sq.Eq{"allocation_id": allocationID}).
			RunWith(tx).
			ExecContext(ctx)
		return err
	})
}

func (s *Allocation) FetchAllocations(ctx context.Context) ([]model.Allocation, error) {
	q := sq.Select("DISTINCT(allocation_id)").
		From("allocations").
		RunWith(s.db)

	rows, err := q.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var allocations []model.Allocation
	for rows.Next() {
		var alloc model.Allocation
		if err := rows.Scan(&alloc.ID); err != nil {
			return nil, err
		}
		allocations = append(allocations, alloc)
	}
	return allocations, nil
}
