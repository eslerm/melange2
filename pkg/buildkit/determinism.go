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
	"maps"
	"slices"

	"github.com/moby/buildkit/client/llb"
)

// SortedEnvOpts returns llb.RunOption entries for environment variables,
// sorted by key for deterministic LLB generation.
//
// CRITICAL: Go map iteration order is random. If we don't sort, the same
// build configuration will produce different LLB digests on different runs,
// breaking caching and reproducibility.
func SortedEnvOpts(env map[string]string) []llb.RunOption {
	if len(env) == 0 {
		return nil
	}

	keys := slices.Sorted(maps.Keys(env))
	opts := make([]llb.RunOption, 0, len(keys))
	for _, k := range keys {
		opts = append(opts, llb.AddEnv(k, env[k]))
	}
	return opts
}

// MergeEnv merges multiple environment maps, with later maps taking precedence.
// Returns a new map without modifying the inputs.
func MergeEnv(envs ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, env := range envs {
		for k, v := range env {
			result[k] = v
		}
	}
	return result
}

// SortedEnvSlice returns environment variables as a sorted slice of "KEY=VALUE" strings.
// This is useful for places that need []string instead of llb.RunOption.
func SortedEnvSlice(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}

	keys := slices.Sorted(maps.Keys(env))
	result := make([]string, 0, len(keys))
	for _, k := range keys {
		result = append(result, k+"="+env[k])
	}
	return result
}
