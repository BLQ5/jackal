/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package mysql

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/ortuman/jackal/model"
)

type Allocations struct {
	*mySQLStorage
}

func newAllocations(db *sql.DB) *Allocations {
	return &Allocations{
		mySQLStorage: newStorage(db),
	}
}

func (s *Allocations) RegisterAllocation(ctx context.Context, allocation *model.Allocation) error {
	_, err := sq.Insert("allocations").
		Columns("allocation_id", "updated_at", "created_at").
		Suffix("ON DUPLICATE KEY UPDATE updated_at = NOW()").
		Values(allocation.ID, nowExpr, nowExpr).
		RunWith(s.db).ExecContext(ctx)
	return err
}

func (s *Allocations) UnregisterAllocation(ctx context.Context, allocationID string) error {
	return s.inTransaction(ctx, func(tx *sql.Tx) error {
		_, err := sq.Delete("presences").
			Where(sq.Eq{"allocation_id": allocationID}).
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return err
		}
		_, err = sq.Delete("resources").
			Where(sq.Eq{"allocation_id": allocationID}).
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return err
		}
		_, err = sq.Delete("allocations").
			Where(sq.Eq{"allocation_id": allocationID}).
			RunWith(tx).
			ExecContext(ctx)
		return err
	})
}

func (s *Allocations) FetchAllocations(ctx context.Context) ([]model.Allocation, error) {
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
