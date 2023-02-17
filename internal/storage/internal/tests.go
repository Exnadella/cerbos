// Copyright 2021-2023 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

//go:build tests
// +build tests

package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cerbos/cerbos/internal/storage"
)

// MutateStoreFn points to a function which mutates the store (ex: add, delete a policy).
type MutateStoreFn func() error

func TestSuiteReloadable(store storage.Store, addFn, deleteFn MutateStoreFn) func(*testing.T) {
	//nolint:thelper
	return func(t *testing.T) {
		r, ok := store.(storage.Reloadable)
		require.True(t, ok, "Store is not reloadable")

		policies, err := store.ListPolicyIDs(context.Background(), false)
		require.NoError(t, err)
		require.Len(t, policies, 0)

		err = addFn()
		require.NoError(t, err)

		policies, err = store.ListPolicyIDs(context.Background(), false)
		require.NoError(t, err)
		require.Len(t, policies, 0)

		err = r.Reload(context.Background())
		require.NoError(t, err)

		policies, err = store.ListPolicyIDs(context.Background(), false)
		require.NoError(t, err)
		require.NotZero(t, len(policies))

		err = deleteFn()
		require.NoError(t, err)

		err = r.Reload(context.Background())
		require.NoError(t, err)

		policies, err = store.ListPolicyIDs(context.Background(), false)
		require.NoError(t, err)
		require.Len(t, policies, 0)
	}
}
