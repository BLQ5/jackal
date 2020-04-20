/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package pgsql

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/ortuman/jackal/model"
	"github.com/ortuman/jackal/xmpp/jid"
)

type Resources struct {
	*pgSQLStorage
}

func newResources(db *sql.DB) *Resources {
	return &Resources{
		pgSQLStorage: newStorage(db),
	}
}

func (r *Resources) UpsertResource(ctx context.Context, resource *model.Resource) error {
	_, err := sq.Insert("resources").
		Columns("allocation_id", "username", "domain", "resource", "priority").
		Values(resource.AllocationID, resource.JID.Node(), resource.JID.Domain(), resource.JID.Resource(), resource.Priority).
		Suffix("ON CONFLICT(username, domain, resource) DO UPDATE SET allocation_id = $1, priority = $5").
		RunWith(r.db).
		ExecContext(ctx)
	return err
}

func (r *Resources) DeleteResource(ctx context.Context, username, domain, resource string) error {
	_, err := sq.Delete("resources").
		Where(sq.And{
			sq.Eq{"username": username},
			sq.Eq{"domain": domain},
			sq.Eq{"resource": resource},
		}).
		RunWith(r.db).
		ExecContext(ctx)
	return err
}

func (r *Resources) FetchResources(ctx context.Context, username, domain string) ([]model.Resource, error) {
	rows, err := sq.Select("allocation_id", "username", "domain", "resource", "priority").
		Where(sq.And{
			sq.Eq{"username": username},
			sq.Eq{"domain": domain},
		}).
		RunWith(r.db).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var resources []model.Resource
	for rows.Next() {
		var allocID, jidStr string
		var priority int8

		if err := rows.Scan(&allocID, &jidStr, &priority); err != nil {
			return nil, err
		}
		resJID, _ := jid.NewWithString(jidStr, true)
		resources = append(resources, model.Resource{
			AllocationID: allocID,
			JID:          resJID,
			Priority:     priority,
		})
	}
	return resources, nil
}
