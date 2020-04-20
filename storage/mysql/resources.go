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
	"github.com/ortuman/jackal/xmpp/jid"
)

type Resources struct {
	*mySQLStorage
}

func newResources(db *sql.DB) *Resources {
	return &Resources{
		mySQLStorage: newStorage(db),
	}
}

func (r *Resources) UpsertResource(ctx context.Context, resource *model.Resource) error {
	_, err := sq.Insert("resources").
		Columns("allocation_id", "username", "domain", "resource", "priority").
		Values(resource.AllocationID, resource.JID.Node(), resource.JID.Domain(), resource.JID.Resource(), resource.Priority).
		Suffix("ON DUPLICATE KEY UPDATE allocation_id = ?, priority = ?, updated_at = NOW()", resource.AllocationID, resource.Priority).
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
		From("resources").
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
		var allocID, username, domain, resource string
		var priority int8

		if err := rows.Scan(&allocID, &username, &domain, &resource, &priority); err != nil {
			return nil, err
		}
		resJID, _ := jid.New(username, domain, resource, true)
		resources = append(resources, model.Resource{
			AllocationID: allocID,
			JID:          resJID,
			Priority:     priority,
		})
	}
	return resources, nil
}
