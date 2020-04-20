/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package model

import (
	"bytes"
	"testing"

	capsmodel "github.com/ortuman/jackal/model/capabilities"
	"github.com/ortuman/jackal/xmpp"
	"github.com/ortuman/jackal/xmpp/jid"
	"github.com/stretchr/testify/require"
)

func TestExtPresence(t *testing.T) {
	j1, _ := jid.NewWithString("ortuman@jackal.im", true)

	var p1, p2 ExtPresence
	p1 = ExtPresence{
		Presence: xmpp.NewPresence(j1, j1, xmpp.AvailableType),
	}

	buf := new(bytes.Buffer)
	require.Nil(t, p1.ToBytes(buf))
	require.Nil(t, p2.FromBytes(buf))
	require.Equal(t, p1, p2)

	var p3, p4 ExtPresence
	p3 = ExtPresence{
		Presence: xmpp.NewPresence(j1, j1, xmpp.AvailableType),
		Caps: &capsmodel.Capabilities{
			Node: "http://jackal.im",
			Ver:  "v1234",
		},
	}
	buf = new(bytes.Buffer)
	require.Nil(t, p3.ToBytes(buf))
	require.Nil(t, p4.FromBytes(buf))
	require.Equal(t, p3, p4)
}
