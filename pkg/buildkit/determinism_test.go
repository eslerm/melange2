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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortedEnvOpts(t *testing.T) {
	env := map[string]string{
		"Z_VAR": "z-value",
		"A_VAR": "a-value",
		"M_VAR": "m-value",
	}

	// Run 100 times to catch non-determinism
	var results []int
	for i := 0; i < 100; i++ {
		opts := SortedEnvOpts(env)
		results = append(results, len(opts))
	}

	// All results should have same length
	for i := 1; i < len(results); i++ {
		require.Equal(t, results[0], results[i])
	}

	// Should have 3 options
	require.Equal(t, 3, results[0])
}

func TestSortedEnvOptsEmpty(t *testing.T) {
	opts := SortedEnvOpts(nil)
	require.Nil(t, opts)

	opts = SortedEnvOpts(map[string]string{})
	require.Nil(t, opts)
}

func TestMergeEnv(t *testing.T) {
	env1 := map[string]string{
		"A": "1",
		"B": "2",
	}
	env2 := map[string]string{
		"B": "override",
		"C": "3",
	}

	merged := MergeEnv(env1, env2)

	require.Equal(t, "1", merged["A"])
	require.Equal(t, "override", merged["B"]) // env2 takes precedence
	require.Equal(t, "3", merged["C"])

	// Original maps should not be modified
	require.Equal(t, "2", env1["B"])
}

func TestMergeEnvEmpty(t *testing.T) {
	merged := MergeEnv()
	require.NotNil(t, merged)
	require.Empty(t, merged)

	merged = MergeEnv(nil, nil)
	require.NotNil(t, merged)
	require.Empty(t, merged)
}

func TestSortedEnvSlice(t *testing.T) {
	env := map[string]string{
		"Z_VAR": "z-value",
		"A_VAR": "a-value",
		"M_VAR": "m-value",
	}

	// Run 100 times to verify determinism
	var firstResult []string
	for i := 0; i < 100; i++ {
		result := SortedEnvSlice(env)
		if i == 0 {
			firstResult = result
		} else {
			require.Equal(t, firstResult, result, "iteration %d should match", i)
		}
	}

	// Should be sorted alphabetically
	require.Equal(t, []string{
		"A_VAR=a-value",
		"M_VAR=m-value",
		"Z_VAR=z-value",
	}, firstResult)
}

func TestSortedEnvSliceEmpty(t *testing.T) {
	result := SortedEnvSlice(nil)
	require.Nil(t, result)

	result = SortedEnvSlice(map[string]string{})
	require.Nil(t, result)
}
