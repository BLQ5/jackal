/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package memorystorage

import (
	"context"
	"testing"

	"github.com/ortuman/jackal/model"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_Allocations(t *testing.T) {
	s := NewAllocations()

	_ = s.RegisterAllocation(context.Background(), &model.Allocation{ID: "a1234"})
	_ = s.RegisterAllocation(context.Background(), &model.Allocation{ID: "b1234"})

	allocs, _ := s.FetchAllocations(context.Background())
	require.Len(t, allocs, 2)

	_ = s.UnregisterAllocation(context.Background(), "a1234")
	_ = s.UnregisterAllocation(context.Background(), "b1234")
}
