/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package model

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllocation(t *testing.T) {

	var a1, a2 Allocation
	a1 = Allocation{
		ID: "a1234",
	}

	buf := new(bytes.Buffer)
	require.Nil(t, a1.ToBytes(buf))
	require.Nil(t, a2.FromBytes(buf))
	require.Equal(t, a1, a2)
}
