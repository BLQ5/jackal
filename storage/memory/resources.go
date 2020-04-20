/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package memorystorage

import (
	"context"

	"github.com/ortuman/jackal/model"
)

type Resources struct {
	*memoryStorage
}

func NewResources() *Resources {
	return &Resources{memoryStorage: newStorage()}
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
