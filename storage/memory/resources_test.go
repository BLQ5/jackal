/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package memorystorage

import (
	"context"
	"testing"

	"github.com/ortuman/jackal/model"
	"github.com/ortuman/jackal/xmpp/jid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_Resources(t *testing.T) {
	j1, _ := jid.NewWithString("ortuman@jackal.im/yard", true)
	j2, _ := jid.NewWithString("ortuman@jackal.im/chamber", true)

	s := NewResources()

	_ = s.UpsertResource(context.Background(), &model.Resource{
		AllocationID: "a1234",
		JID:          j1,
		Priority:     1,
	})
	_ = s.UpsertResource(context.Background(), &model.Resource{
		AllocationID: "a5678",
		JID:          j2,
		Priority:     8,
	})

	resources, _ := s.FetchResources(context.Background(), "ortuman", "jackal.im")
	require.Len(t, resources, 2)

	_ = s.DeleteResource(context.Background(), "ortuman", "jackal.im", "yard")

	resources, _ = s.FetchResources(context.Background(), "ortuman", "jackal.im")
	require.Len(t, resources, 1)

	require.Equal(t, "a5678", resources[0].AllocationID)
	require.Equal(t, "ortuman@jackal.im/chamber", resources[0].JID.String())
	require.Equal(t, int8(8), resources[0].Priority)
}
