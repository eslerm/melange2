// Copyright 2024 Chainguard, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package buildkit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClientDefaultAddr(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	bk := startBuildKitContainer(t, ctx)

	// Test with explicit address
	c, err := New(ctx, bk.Addr)
	require.NoError(t, err)
	require.Equal(t, bk.Addr, c.Addr())

	require.NoError(t, c.Ping(ctx))
	require.NoError(t, c.Close())
}

func TestNewClientConnectionError(t *testing.T) {
	ctx := context.Background()

	// Try to connect to a non-existent address
	// Note: The BuildKit client connection is lazy, so New() might succeed
	// but Ping() will fail
	c, err := New(ctx, "tcp://localhost:59999")
	if err != nil {
		// Connection failed immediately (expected on some systems)
		require.Contains(t, err.Error(), "connecting to buildkit")
		return
	}
	defer c.Close()

	// Connection was lazy - Ping should fail
	err = c.Ping(ctx)
	require.Error(t, err)
}

func TestClientPing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	bk := startBuildKitContainer(t, ctx)

	c, err := New(ctx, bk.Addr)
	require.NoError(t, err)
	defer c.Close()

	// Should succeed
	require.NoError(t, c.Ping(ctx))

	// Underlying client should be accessible
	require.NotNil(t, c.Client())
}
