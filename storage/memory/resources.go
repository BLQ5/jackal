/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package memorystorage

import (
	"context"
	"strings"

	"github.com/ortuman/jackal/model"
	"github.com/ortuman/jackal/model/serializer"
)

type Resources struct {
	*memoryStorage
}

func NewResources() *Resources {
	return &Resources{memoryStorage: newStorage()}
}

func (m *Resources) UpsertResource(_ context.Context, resource *model.Resource) error {
	return m.saveEntity(resourceKey(resource.JID.Node(), resource.JID.Domain(), resource.JID.Resource()), resource)
}

func (m *Resources) DeleteResource(_ context.Context, username, domain, resource string) error {
	return m.deleteKey(resourceKey(username, domain, resource))
}

func (m *Resources) FetchResources(_ context.Context, username, domain string) ([]model.Resource, error) {
	var resources []model.Resource
	if err := m.inReadLock(func() error {
		for k, b := range m.b {
			if !strings.HasPrefix(k, "resources:"+username+":"+domain) {
				continue
			}
			var res model.Resource
			if err := serializer.Deserialize(b, &res); err != nil {
				return err
			}
			resources = append(resources, res)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return resources, nil
}

func resourceKey(username, domain, resource string) string {
	return "resources:" + username + ":" + domain + ":" + resource
}
