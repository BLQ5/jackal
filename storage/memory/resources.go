/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package memorystorage

import (
	"github.com/ortuman/jackal/model"
)

type Resources struct {
	*memoryStorage
}

func NewResources() *Resources {
	return &Resources{memoryStorage: newStorage()}
}

func (r *Resources) UpsertResource(resource *model.Resource, allocationID string) error {
	panic("implement me!")
}

func (r *Resources) FetchResources(username, domain string) ([]model.Resource, error) {
	panic("implement me!")
}
