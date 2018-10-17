# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [0.4.0] - 2018-10-17
### Fixed
- Fix subrouters appending middlewares to the parent's middleware slice.

### Removed
- `(*Router).SetCtxKey` in order to set the key only once, when creating a new router.

## [0.3.2] - 2018-10-03
### Fixed
- Copy placeholder when creating a subrouter.

## [0.3.1] - 2018-10-03
### Fixed
- Return an empty map when URL parameters don't exist.

## [0.3.0] - 2018-09-25
### Added
- Test routine on AppVeyor.
- Third-party dependency for the radix tree.

### Changed
- Update `doc.go` with a simpler text.
- Merge `Mux` and `Router` types into a single type.
- Force setting a context key for URL parameters when creating a router.
- Change constructors' names.

### Removed
- `internal` package.

## [0.2.0] - 2018-08-21
### Added
- Router (prefixed mux).

## [0.1.1] - 2018-08-21
### Added
- New test cases with empty paths.

### Changed
- Empty path defaults to single slash.

### Fixed
- Infinite loop when path is an empty string.

## 0.1.0 - 2018-07-31
### Added
- This changelog file.
- README file.
- MIT License.
- Travis CI configuration file.
- Makefile.
- Git ignore file.
- EditorConfig file.
- This package's source code, including examples and tests.

[0.4.0]: https://github.com/gbrlsnchs/mux/compare/v0.3.2...v0.4.0
[0.3.2]: https://github.com/gbrlsnchs/mux/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/gbrlsnchs/mux/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/gbrlsnchs/mux/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/gbrlsnchs/mux/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/gbrlsnchs/mux/compare/v0.1.0...v0.1.1