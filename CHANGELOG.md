# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
