/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package storage

import (
	"context"

	"github.com/ortuman/jackal/model"
)

type Resources interface {
	UpsertResource(ctx context.Context, resource *model.Resource) error
	DeleteResource(ctx context.Context, username, domain, resource string) error

	FetchResources(ctx context.Context, username, domain string) ([]model.Resource, error)
}
