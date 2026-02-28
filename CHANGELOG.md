# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Added

- Statistics page showcasing Evaporation and Efficiency from past recipes

### Changed

- Added proper SQL migrations
- Added seed data for statistics from past recipes
- Update dependencies
  - `github.com/knadh/koanf/providers/file` v1.2.0 -> v1.2.1
	- `github.com/knadh/koanf/v2` v2.2.2 -> v2.3.2
	- `github.com/labstack/echo/v4` v4.13.4 -> v4.15.1
	- `github.com/mattn/go-sqlite3` v1.14.28 -> v1.14.34
	- `github.com/stretchr/testify` v1.10.0 -> v1.11.1

## [2.0.5] - 2025-07-19

### Changed
- Updated to go 1.24
- Update dependencies
  - `github.com/knadh/koanf/parsers/yaml` v0.1.0 -> v1.1.0
  - `github.com/knadh/koanf/providers/env` v0.1.0 -> v1.1.0
  - `github.com/knadh/koanf/providers/file` v0.1.0 -> v1.2.0
  - `github.com/knadh/koanf/v2` v2.0.1 -> v2.2.2
  - `github.com/labstack/echo/v4` v.11.1 -> v4.13.4
  - `github.com/mattn/go-sqlite3` v.1.14.22 -> v1.14.28
  - `github.com/rs/zerolog` v1.30.0 -> v1.34.0
  - `github.com/stretchr/testify` v.1.8.4 -> v1.10.0
- Update actions in pipeline to latest
- Update gotify recommendation

## [2.0.4] - 2025-07-19

### Changed

- New logo that is way cooler

## [2.0.3] - 2024-06-08

### Changed

- ARM Base image is now `arm64v8/debian:bookworm-slim`. This fixes the execution issues but makes the image considerably more heavy (~100MB)

## [2.0.2] - 2024-06-08

- Fix CGO compilation in ARM targets

## [2.0.1] - 2024-06-08

- Compile binaries with CGO=1 to enable sqlite 

## [2.0.0] - 2024-06-08

### Added
- Option in frontend to select summary format
- Option to add notes in each hop
- Persistent storage option based on `sqlite` to enable application restarts without data loss
- New Dry Hopping front end that makes sense

## [1.0.1] - 2024-02-02

### Added
- Fixed status bug where recipes that are just starting had a "Unknown" status, and could not be continued

## [1.0.0] - 2023-11-21

### Added
- Recipes page with a list of all active recipes and possibility to continue a recipe.
- Extended fermentation with notification setting pages. Now, the user can choose which notifications to receive and when.
- Added dry hop , secondary fermentation and bottling to the process.

### Changed
- Fixed efficiency calculation for mash.


## [0.3.0] - 2023-10-01

### Added
- Extended Mash View: Now, an overview of missings rasts and their duration can be seen from any rast.

### Changed
- Improved color scale for estimating EBC. The scale now contains better colors for the higher values and covers a wider range of values. The method is now also more robust, avoiding panics when the EBC is out of range.

## [0.2.0] - 2023-09-25

### Added
- Configuration trough YAML file and/or environment variables.
- Notifications for timers via Gotify

## [0.1.1] - 2023-09-09
### Changed

- Changed "Läuterruhe" from 1 to 15 minutes.

## [0.1.0] - 2023-09-09

Initial release.
