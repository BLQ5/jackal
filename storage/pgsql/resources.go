/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package pgsql

import (
	"context"
	"database/sql"

	"github.com/ortuman/jackal/model"
)

type Resources struct {
	*pgSQLStorage
}

func newResources(db *sql.DB) *Resources {
	return &Resources{
		pgSQLStorage: newStorage(db),
	}
}

func (r *Resources) UpsertResource(ctx context.Context, resource *model.Resource, allocationID string) error {
	panic("implement me!")
}

func (r *Resources) DeleteResource(ctx context.Context, username, domain, resource string) error {
	panic("implement me!")
}

func (r *Resources) FetchResources(ctx context.Context, username, domain string) ([]model.Resource, error) {
	panic("implement me!")
}
