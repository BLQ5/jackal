/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package mysql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ortuman/jackal/model"
	"github.com/stretchr/testify/require"
)

func TestMySQLStorage_AllocationsRegister(t *testing.T) {
	s, mock := newAllocationsMock()

	mock.ExpectExec("INSERT INTO allocations (.+) VALUES(.+) ON DUPLICATE KEY UPDATE updated_at = NOW\\(\\)").
		WithArgs("a1234").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.RegisterAllocation(context.Background(), &model.Allocation{
		ID: "a1234",
	})
	require.Nil(t, err)

	require.Nil(t, mock.ExpectationsWereMet())
}

func TestMySQLStorage_AllocationsUnregister(t *testing.T) {
	s, mock := newAllocationsMock()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM presences WHERE allocation_id = \\?").
		WithArgs("a1234").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM resources WHERE allocation_id = \\?").
		WithArgs("a1234").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM allocations WHERE allocation_id = \\?").
		WithArgs("a1234").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := s.UnregisterAllocation(context.Background(), "a1234")
	require.Nil(t, err)

	require.Nil(t, mock.ExpectationsWereMet())
}

func TestMySQLStorage_AllocationsFetchAllocations(t *testing.T) {
	rows := sqlmock.NewRows([]string{"allocation_id"})
	rows.AddRow(`a1234`)
	rows.AddRow(`b1234`)

	s, mock := newAllocationsMock()

	mock.ExpectQuery("SELECT DISTINCT\\(allocation_id\\) FROM allocations").
		WillReturnRows(rows)

	allocs, err := s.FetchAllocations(context.Background())
	require.Nil(t, err)
	require.Len(t, allocs, 2)

	require.Nil(t, mock.ExpectationsWereMet())
}

func newAllocationsMock() (*Allocations, sqlmock.Sqlmock) {
	s, sqlMock := newStorageMock()
	return &Allocations{
		mySQLStorage: s,
	}, sqlMock
}
