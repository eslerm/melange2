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

package storage

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/googleapi"
)

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		retryable bool
	}{
		{
			name:      "nil error",
			err:       nil,
			retryable: false,
		},
		{
			name:      "429 too many requests",
			err:       &googleapi.Error{Code: 429},
			retryable: true,
		},
		{
			name:      "500 internal server error",
			err:       &googleapi.Error{Code: 500},
			retryable: true,
		},
		{
			name:      "502 bad gateway",
			err:       &googleapi.Error{Code: 502},
			retryable: true,
		},
		{
			name:      "503 service unavailable",
			err:       &googleapi.Error{Code: 503},
			retryable: true,
		},
		{
			name:      "504 gateway timeout",
			err:       &googleapi.Error{Code: 504},
			retryable: true,
		},
		{
			name:      "400 bad request - not retryable",
			err:       &googleapi.Error{Code: 400},
			retryable: false,
		},
		{
			name:      "401 unauthorized - not retryable",
			err:       &googleapi.Error{Code: 401},
			retryable: false,
		},
		{
			name:      "403 forbidden - not retryable",
			err:       &googleapi.Error{Code: 403},
			retryable: false,
		},
		{
			name:      "404 not found - not retryable",
			err:       &googleapi.Error{Code: 404},
			retryable: false,
		},
		{
			name:      "context deadline exceeded",
			err:       context.DeadlineExceeded,
			retryable: true,
		},
		{
			name:      "context canceled - not retryable",
			err:       context.Canceled,
			retryable: false,
		},
		{
			name:      "connection reset",
			err:       errors.New("connection reset by peer"),
			retryable: true,
		},
		{
			name:      "connection refused",
			err:       errors.New("connection refused"),
			retryable: true,
		},
		{
			name:      "timeout error",
			err:       errors.New("i/o timeout"),
			retryable: true,
		},
		{
			name:      "temporary failure",
			err:       errors.New("temporary failure in name resolution"),
			retryable: true,
		},
		{
			name:      "wrapped googleapi error",
			err:       errors.Join(errors.New("wrapper"), &googleapi.Error{Code: 503}),
			retryable: true,
		},
		{
			name:      "regular error - not retryable",
			err:       errors.New("some random error"),
			retryable: false,
		},
		{
			name:      "net.OpError with timeout",
			err:       &net.OpError{Err: errors.New("i/o timeout")},
			retryable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRetryableError(tt.err)
			assert.Equal(t, tt.retryable, result)
		})
	}
}

func TestGCSStorageDefaults(t *testing.T) {
	// Test that default values are correct
	assert.Equal(t, 50, DefaultMaxConcurrentUploads)
	assert.Equal(t, 5, DefaultMaxRetries)
	assert.Equal(t, 100*time.Millisecond, DefaultInitialBackoff)
	assert.Equal(t, 30*time.Second, DefaultMaxBackoff)
}

func TestGCSOptions(t *testing.T) {
	t.Run("WithMaxConcurrentUploads", func(t *testing.T) {
		s := &GCSStorage{}
		opt := WithMaxConcurrentUploads(100)
		opt(s)
		assert.Equal(t, 100, s.maxConcurrentUploads)
		assert.Equal(t, 100, cap(s.uploadSem))
	})

	t.Run("WithRetryConfig", func(t *testing.T) {
		s := &GCSStorage{}
		opt := WithRetryConfig(10, 200*time.Millisecond, 60*time.Second)
		opt(s)
		assert.Equal(t, 10, s.maxRetries)
		assert.Equal(t, 200*time.Millisecond, s.initialBackoff)
		assert.Equal(t, 60*time.Second, s.maxBackoff)
	})
}
