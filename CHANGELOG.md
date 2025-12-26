# Changelog

All notable changes to melange2 will be documented in this file.

## [Unreleased]

### Added
- BuildKit as the sole execution backend for `build` and `rebuild` commands
- Real-time progress display with step numbers, names, and timing
- Cache hit/miss information in build output
- Stdout/stderr streaming from build steps (with `--debug`)
- Build summary showing total steps, cached steps, and duration
- `--buildkit-addr` flag to specify BuildKit daemon address

### Changed
- Removed bubblewrap, Docker, and QEMU runners from build command
- Simplified architecture with single execution backend

### Known Limitations
- `test` command still uses legacy runner system
- Requires external BuildKit daemon

## Fork History

This project is forked from [chainguard-dev/melange](https://github.com/chainguard-dev/melange).
See the upstream project for historical changelog entries.
