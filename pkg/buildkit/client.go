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

// Package buildkit provides BuildKit integration for melange builds.
// It replaces the previous runner implementations (bubblewrap, docker, qemu)
// with a single BuildKit-based execution backend.
package buildkit

import (
	"context"
	"fmt"

	"github.com/moby/buildkit/client"
)

const (
	// DefaultAddr is the default BuildKit daemon address.
	DefaultAddr = "tcp://localhost:1234"
)

// Client wraps the BuildKit client with melange-specific functionality.
type Client struct {
	bk   *client.Client
	addr string
}

// New creates a new BuildKit client connected to the specified address.
// If addr is empty, DefaultAddr is used.
func New(ctx context.Context, addr string) (*Client, error) {
	if addr == "" {
		addr = DefaultAddr
	}

	bk, err := client.New(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("connecting to buildkit at %s: %w", addr, err)
	}

	return &Client{
		bk:   bk,
		addr: addr,
	}, nil
}

// Close closes the connection to BuildKit.
func (c *Client) Close() error {
	if c.bk != nil {
		return c.bk.Close()
	}
	return nil
}

// Addr returns the BuildKit daemon address.
func (c *Client) Addr() string {
	return c.addr
}

// Ping verifies the connection to BuildKit is working.
func (c *Client) Ping(ctx context.Context) error {
	workers, err := c.bk.ListWorkers(ctx)
	if err != nil {
		return fmt.Errorf("pinging buildkit: %w", err)
	}
	if len(workers) == 0 {
		return fmt.Errorf("buildkit has no workers")
	}
	return nil
}

// Client returns the underlying BuildKit client for advanced operations.
func (c *Client) Client() *client.Client {
	return c.bk
}
