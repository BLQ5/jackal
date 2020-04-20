/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package model

import (
	"bytes"
	"testing"

	"github.com/ortuman/jackal/xmpp/jid"
	"github.com/stretchr/testify/require"
)

func TestResource(t *testing.T) {
	j, _ := jid.NewWithString("ortuman@jackal.im/yard", true)

	var r1, r2 Resource
	r1 = Resource{
		AllocationID: "a1234",
		JID:          j,
		Priority:     8,
	}

	buf := new(bytes.Buffer)
	require.Nil(t, r1.ToBytes(buf))
	require.Nil(t, r2.FromBytes(buf))
	require.Equal(t, r1, r2)
}
