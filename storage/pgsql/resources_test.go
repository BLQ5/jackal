/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package pgsql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ortuman/jackal/model"
	"github.com/ortuman/jackal/xmpp/jid"
	"github.com/stretchr/testify/require"
)

func TestPgSQLStorage_ResourcesUpsert(t *testing.T) {
	j1, _ := jid.NewWithString("ortuman@jackal.im/yard", true)

	s, mock := newResourcesMock()

	mock.ExpectExec("INSERT INTO resources (.+) VALUES (.+) ON CONFLICT\\(username, domain, resource\\) DO UPDATE SET allocation_id = \\$1, priority = \\$5").
		WithArgs("a1234", "ortuman", "jackal.im", "yard", 8).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.UpsertResource(context.Background(), &model.Resource{
		AllocationID: "a1234",
		JID:          j1,
		Priority:     8,
	})
	require.Nil(t, err)

	require.Nil(t, mock.ExpectationsWereMet())
}

func TestPgSQLStorage_ResourcesDelete(t *testing.T) {
	s, mock := newResourcesMock()

	mock.ExpectExec("DELETE FROM resources WHERE \\(username = \\? AND domain = \\? AND resource = \\?\\)").
		WithArgs("ortuman", "jackal.im", "yard").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.DeleteResource(context.Background(), "ortuman", "jackal.im", "yard")
	require.Nil(t, err)

	require.Nil(t, mock.ExpectationsWereMet())
}

func TestPgSQLStorage_ResourcesFetch(t *testing.T) {
	rows := sqlmock.NewRows([]string{"allocation_id", "username", "domain", "resource", "priority"})
	rows.AddRow(`a1234`, `ortuman`, `jackal.im`, `yard`, 1)
	rows.AddRow(`b1234`, `ortuman`, `jackal.im`, `chamber`, 8)

	s, mock := newResourcesMock()

	mock.ExpectQuery("SELECT allocation_id, username, domain, resource, priority FROM resources WHERE \\(username = \\? AND domain = \\?\\)").
		WithArgs("ortuman", "jackal.im").
		WillReturnRows(rows)

	resources, err := s.FetchResources(context.Background(), "ortuman", "jackal.im")
	require.Nil(t, err)
	require.Len(t, resources, 2)

	require.Nil(t, mock.ExpectationsWereMet())
}

func newResourcesMock() (*Resources, sqlmock.Sqlmock) {
	s, sqlMock := newStorageMock()
	return &Resources{
		pgSQLStorage: s,
	}, sqlMock
}
