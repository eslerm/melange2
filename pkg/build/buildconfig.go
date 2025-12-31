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

package build

import (
	"fmt"
	"os"
	"time"

	apko_types "chainguard.dev/apko/pkg/build/types"
	"chainguard.dev/apko/pkg/options"

	"github.com/dlorenc/melange2/pkg/config"
)

// BuildConfig contains all immutable configuration for a build.
// This struct is the single source of truth for build parameters and can be
// populated from CLI flags, remote build specs, or programmatically.
type BuildConfig struct {
	// ConfigFile is the path to the build configuration file (e.g., "crane.yaml").
	ConfigFile string

	// Configuration is the parsed melange configuration.
	// If set, ConfigFile is used only for reporting purposes.
	Configuration *config.Configuration

	// ConfigFileRepositoryURL is the URL of the git repository where the build
	// configuration file is stored (e.g., "https://github.com/wolfi-dev/os").
	ConfigFileRepositoryURL string

	// ConfigFileRepositoryCommit is the commit hash of the git repository
	// corresponding to the current state of the build configuration file.
	ConfigFileRepositoryCommit string

	// ConfigFileLicense is the SPDX license string for the build configuration file.
	ConfigFileLicense string

	// SourceDateEpoch is the timestamp used for reproducible builds.
	SourceDateEpoch time.Time

	// WorkspaceDir is the directory used for the build workspace at /home/build.
	WorkspaceDir string

	// WorkspaceIgnore is the file containing ignore rules for the workspace.
	WorkspaceIgnore string

	// PipelineDirs are ordered directories where to find 'uses' pipelines.
	PipelineDirs []string

	// SourceDir is the directory containing source files for the build.
	SourceDir string

	// SigningKey is the path to the key used for signing packages.
	SigningKey string

	// SigningPassphrase is the passphrase for the signing key.
	SigningPassphrase string

	// Namespace is the namespace used in package URLs in SBOM.
	Namespace string

	// GenerateIndex indicates whether to generate APKINDEX.tar.gz.
	GenerateIndex bool

	// EmptyWorkspace indicates whether the build workspace should be empty.
	EmptyWorkspace bool

	// OutDir is the directory where packages will be output.
	OutDir string

	// Arch is the target architecture for the build.
	Arch apko_types.Architecture

	// Libc is the libc flavor override (e.g., "gnu", "musl").
	Libc string

	// ExtraKeys are additional keys to include in the build environment keyring.
	ExtraKeys []string

	// ExtraRepos are additional repositories to include in the build environment.
	ExtraRepos []string

	// ExtraPackages are extra packages to install for the build environment.
	ExtraPackages []string

	// DependencyLog is the filename for dependency logging.
	DependencyLog string

	// CreateBuildLog indicates whether to generate a package.log file.
	CreateBuildLog bool

	// PersistLintResults indicates whether to persist lint results to JSON files.
	PersistLintResults bool

	// CacheDir is the directory used for cached inputs.
	CacheDir string

	// ApkCacheDir is the directory used for cached apk packages.
	ApkCacheDir string

	// StripOriginName determines whether origin names should be stripped.
	StripOriginName bool

	// EnvFile is the environment file for preloading build environment variables.
	EnvFile string

	// VarsFile is the variables file for build configuration variables.
	VarsFile string

	// BuildKitAddr is the BuildKit daemon address.
	BuildKitAddr string

	// Debug enables debug logging of build pipelines.
	Debug bool

	// Remove indicates whether to clean up intermediate artifacts.
	Remove bool

	// CacheRegistry is the registry URL for BuildKit cache.
	CacheRegistry string

	// CacheMode is the cache export mode ("min" or "max").
	CacheMode string

	// ApkoRegistry is the registry URL for caching apko base images.
	ApkoRegistry string

	// ApkoRegistryInsecure allows insecure connection to ApkoRegistry.
	ApkoRegistryInsecure bool

	// LintRequire are linter checks that must pass.
	LintRequire []string

	// LintWarn are linter checks that generate warnings.
	LintWarn []string

	// Auth contains authentication for package repositories.
	Auth map[string]options.Auth

	// IgnoreSignatures indicates whether to ignore repository signature verification.
	IgnoreSignatures bool

	// EnabledBuildOptions are build options to apply to the configuration.
	EnabledBuildOptions []string

	// MaxLayers controls the maximum number of layers for the build environment.
	MaxLayers int

	// ExportOnFailure specifies how to export the build environment on failure.
	ExportOnFailure string

	// ExportRef is the path or image reference for debug image export.
	ExportRef string

	// GenerateProvenance indicates whether to generate SLSA provenance.
	GenerateProvenance bool
}

// NewBuildConfig creates a new BuildConfig with sensible defaults.
func NewBuildConfig() *BuildConfig {
	return &BuildConfig{
		WorkspaceIgnore: ".melangeignore",
		OutDir:          ".",
		CacheDir:        "./melange-cache/",
		Remove:          true,
		MaxLayers:       50,
	}
}

// ToOptions converts BuildConfig to a slice of build.Option for backward compatibility.
// This allows gradual migration from the Option pattern to direct BuildConfig usage.
func (c *BuildConfig) ToOptions() []Option {
	opts := []Option{}

	if c.ConfigFile != "" {
		opts = append(opts, WithConfig(c.ConfigFile))
	}
	if c.Configuration != nil {
		opts = append(opts, WithConfiguration(c.Configuration, c.ConfigFile))
	}
	if c.ConfigFileRepositoryURL != "" {
		opts = append(opts, WithConfigFileRepositoryURL(c.ConfigFileRepositoryURL))
	}
	if c.ConfigFileRepositoryCommit != "" {
		opts = append(opts, WithConfigFileRepositoryCommit(c.ConfigFileRepositoryCommit))
	}
	if c.ConfigFileLicense != "" {
		opts = append(opts, WithConfigFileLicense(c.ConfigFileLicense))
	}
	if !c.SourceDateEpoch.IsZero() {
		opts = append(opts, WithBuildDate(c.SourceDateEpoch.Format(time.RFC3339)))
	}
	if c.WorkspaceDir != "" {
		opts = append(opts, WithWorkspaceDir(c.WorkspaceDir))
	}
	if c.WorkspaceIgnore != "" && c.WorkspaceIgnore != ".melangeignore" {
		opts = append(opts, WithWorkspaceIgnore(c.WorkspaceIgnore))
	}
	for _, dir := range c.PipelineDirs {
		opts = append(opts, WithPipelineDir(dir))
	}
	if c.SourceDir != "" {
		opts = append(opts, WithSourceDir(c.SourceDir))
	}
	if c.SigningKey != "" {
		opts = append(opts, WithSigningKey(c.SigningKey))
	}
	if c.Namespace != "" {
		opts = append(opts, WithNamespace(c.Namespace))
	}
	opts = append(opts, WithGenerateIndex(c.GenerateIndex))
	opts = append(opts, WithEmptyWorkspace(c.EmptyWorkspace))
	if c.OutDir != "" && c.OutDir != "." {
		opts = append(opts, WithOutDir(c.OutDir))
	}
	if c.Arch != "" {
		opts = append(opts, WithArch(c.Arch))
	}
	if c.Libc != "" {
		opts = append(opts, WithLibcFlavorOverride(c.Libc))
	}
	if len(c.ExtraKeys) > 0 {
		opts = append(opts, WithExtraKeys(c.ExtraKeys))
	}
	if len(c.ExtraRepos) > 0 {
		opts = append(opts, WithExtraRepos(c.ExtraRepos))
	}
	if len(c.ExtraPackages) > 0 {
		opts = append(opts, WithExtraPackages(c.ExtraPackages))
	}
	if c.DependencyLog != "" {
		opts = append(opts, WithDependencyLog(c.DependencyLog))
	}
	opts = append(opts, WithCreateBuildLog(c.CreateBuildLog))
	opts = append(opts, WithPersistLintResults(c.PersistLintResults))
	if c.CacheDir != "" && c.CacheDir != "./melange-cache/" {
		opts = append(opts, WithCacheDir(c.CacheDir))
	}
	if c.ApkCacheDir != "" {
		opts = append(opts, WithPackageCacheDir(c.ApkCacheDir))
	}
	opts = append(opts, WithStripOriginName(c.StripOriginName))
	if c.EnvFile != "" {
		opts = append(opts, WithEnvFile(c.EnvFile))
	}
	if c.VarsFile != "" {
		opts = append(opts, WithVarsFile(c.VarsFile))
	}
	if c.BuildKitAddr != "" {
		opts = append(opts, WithBuildKitAddr(c.BuildKitAddr))
	}
	opts = append(opts, WithDebug(c.Debug))
	opts = append(opts, WithRemove(c.Remove))
	if c.CacheRegistry != "" {
		opts = append(opts, WithCacheRegistry(c.CacheRegistry))
	}
	if c.CacheMode != "" {
		opts = append(opts, WithCacheMode(c.CacheMode))
	}
	if c.ApkoRegistry != "" {
		opts = append(opts, WithApkoRegistry(c.ApkoRegistry))
	}
	opts = append(opts, WithApkoRegistryInsecure(c.ApkoRegistryInsecure))
	if len(c.LintRequire) > 0 {
		opts = append(opts, WithLintRequire(c.LintRequire))
	}
	if len(c.LintWarn) > 0 {
		opts = append(opts, WithLintWarn(c.LintWarn))
	}
	for domain, auth := range c.Auth {
		opts = append(opts, WithAuth(domain, auth.User, auth.Pass))
	}
	opts = append(opts, WithIgnoreSignatures(c.IgnoreSignatures))
	if len(c.EnabledBuildOptions) > 0 {
		opts = append(opts, WithEnabledBuildOptions(c.EnabledBuildOptions))
	}
	if c.MaxLayers != 0 && c.MaxLayers != 50 {
		opts = append(opts, WithMaxLayers(c.MaxLayers))
	}
	if c.ExportOnFailure != "" {
		opts = append(opts, WithExportOnFailure(c.ExportOnFailure, c.ExportRef))
	}
	opts = append(opts, WithGenerateProvenance(c.GenerateProvenance))

	return opts
}

// Validate checks that required fields are set and returns an error if not.
func (c *BuildConfig) Validate() error {
	if c.ConfigFile == "" && c.Configuration == nil {
		return fmt.Errorf("either ConfigFile or Configuration must be set")
	}
	if c.ConfigFileRepositoryURL == "" {
		return fmt.Errorf("ConfigFileRepositoryURL is required")
	}
	if c.ConfigFileRepositoryCommit == "" {
		return fmt.Errorf("ConfigFileRepositoryCommit is required")
	}
	if c.SigningKey != "" {
		if _, err := os.Stat(c.SigningKey); err != nil {
			return fmt.Errorf("signing key not found: %w", err)
		}
	}
	return nil
}

// Clone returns a deep copy of the BuildConfig.
func (c *BuildConfig) Clone() *BuildConfig {
	clone := *c
	// Deep copy slices
	if c.PipelineDirs != nil {
		clone.PipelineDirs = make([]string, len(c.PipelineDirs))
		copy(clone.PipelineDirs, c.PipelineDirs)
	}
	if c.ExtraKeys != nil {
		clone.ExtraKeys = make([]string, len(c.ExtraKeys))
		copy(clone.ExtraKeys, c.ExtraKeys)
	}
	if c.ExtraRepos != nil {
		clone.ExtraRepos = make([]string, len(c.ExtraRepos))
		copy(clone.ExtraRepos, c.ExtraRepos)
	}
	if c.ExtraPackages != nil {
		clone.ExtraPackages = make([]string, len(c.ExtraPackages))
		copy(clone.ExtraPackages, c.ExtraPackages)
	}
	if c.LintRequire != nil {
		clone.LintRequire = make([]string, len(c.LintRequire))
		copy(clone.LintRequire, c.LintRequire)
	}
	if c.LintWarn != nil {
		clone.LintWarn = make([]string, len(c.LintWarn))
		copy(clone.LintWarn, c.LintWarn)
	}
	if c.EnabledBuildOptions != nil {
		clone.EnabledBuildOptions = make([]string, len(c.EnabledBuildOptions))
		copy(clone.EnabledBuildOptions, c.EnabledBuildOptions)
	}
	if c.Auth != nil {
		clone.Auth = make(map[string]options.Auth)
		for k, v := range c.Auth {
			clone.Auth[k] = v
		}
	}
	return &clone
}

// RemoteBuildParams contains parameters for creating a BuildConfig for remote builds.
// This avoids circular dependencies between build and service packages.
type RemoteBuildParams struct {
	ConfigPath   string
	PipelineDir  string
	SourceDir    string
	OutputDir    string
	CacheDir     string
	BackendAddr  string
	Debug        bool
	JobID        string
	CacheRegistry string
	CacheMode     string
	ApkoRegistry  string
	ApkoRegistryInsecure bool
}

// NewBuildConfigForRemote creates a BuildConfig for remote/service builds.
// This is used by the scheduler to convert remote build requests to BuildConfig.
func NewBuildConfigForRemote(params RemoteBuildParams) *BuildConfig {
	cfg := NewBuildConfig()

	cfg.ConfigFile = params.ConfigPath
	cfg.ConfigFileRepositoryURL = "https://melange-service/inline"
	cfg.ConfigFileRepositoryCommit = "inline-" + params.JobID
	cfg.ConfigFileLicense = "Apache-2.0"

	if params.PipelineDir != "" {
		cfg.PipelineDirs = []string{params.PipelineDir}
	}
	if params.SourceDir != "" {
		cfg.SourceDir = params.SourceDir
	}

	cfg.OutDir = params.OutputDir
	cfg.CacheDir = params.CacheDir
	cfg.BuildKitAddr = params.BackendAddr
	cfg.Debug = params.Debug
	cfg.GenerateIndex = true
	cfg.IgnoreSignatures = true
	cfg.Namespace = "wolfi"

	// Cache configuration
	cfg.CacheRegistry = params.CacheRegistry
	cfg.CacheMode = params.CacheMode
	cfg.ApkoRegistry = params.ApkoRegistry
	cfg.ApkoRegistryInsecure = params.ApkoRegistryInsecure

	// Default repos and keys for Wolfi
	cfg.ExtraRepos = []string{"https://packages.wolfi.dev/os"}
	cfg.ExtraKeys = []string{"https://packages.wolfi.dev/os/wolfi-signing.rsa.pub"}

	return cfg
}
