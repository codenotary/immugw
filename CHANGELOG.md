# CHANGELOG
All notable changes to this project will be documented in this file. This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
<a name="unreleased"></a>
## [Unreleased]


<a name="v1.0.5"></a>
## [v1.0.5] - 2021-08-25
### Changes
- bump immudb version
- add coverall
- remove codecov and fix go.mod
- immugw support immudb 0.9.x family
- upgrading for compatibility with immudb at 888ed37bf6cc
- **pkg/api:** clean obsolete swagger
- **pkg/gw:** handle corrupted data as 409 status conflict error


<a name="v0.9.2-RC1"></a>
## [v0.9.2-RC1] - 2021-04-26
### Changes
- add coverall
- remove codecov and fix go.mod
- immugw support immudb 0.9.x family
- upgrading for compatibility with immudb at 888ed37bf6cc
- **pkg/api:** clean obsolete swagger
- **pkg/gw:** handle corrupted data as 409 status conflict error


<a name="v0.8.1"></a>
## v0.8.1 - 2020-12-05
### Bug Fixes
- **cmd/immugw/command:** fix commandline_test error

### Changes
- remove verbose flag from test command

### Code Refactoring
- fixes signed root compatibility usage
- init immugw

### Features
- auditor verifies root signature
- **cmd/immugw:** add service management command


[Unreleased]: https://github.com/vchain-us/immudb/compare/v1.0.5...HEAD
[v1.0.5]: https://github.com/vchain-us/immudb/compare/v0.9.2-RC1...v1.0.5
[v0.9.2-RC1]: https://github.com/vchain-us/immudb/compare/v0.8.1...v0.9.2-RC1
