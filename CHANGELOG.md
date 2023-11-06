# CHANGELOG
All notable changes to this project will be documented in this file. This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
<a name="unreleased"></a>
## [Unreleased]


<a name="v1.3.1"></a>
## [v1.3.1] - 2023-07-27
### Bug Fixes
- **pkg/gw:** serverinfo endpoint is not bounded to any particular database ([#33](https://github.com/vchain-us/immudb/issues/33))

### Changes
- **README:** update swagger usage in README ([#30](https://github.com/vchain-us/immudb/issues/30))
- **pkg/gw:** support multiple kv pairs in set endpoint ([#34](https://github.com/vchain-us/immudb/issues/34))


<a name="v1.3.0"></a>
## [v1.3.0] - 2022-12-07

<a name="v1.3.0-RC1"></a>
## [v1.3.0-RC1] - 2022-12-05
### Bug Fixes
- **ci:** secure docker login command for password ([#24](https://github.com/vchain-us/immudb/issues/24))


<a name="v1.2.0"></a>
## [v1.2.0] - 2022-10-17
### Bug Fixes
- **ci:** update docker credentials

### Changes
- bump immudb version
- bump immudb version
- **ci:** update docker image and build prod images
- **pkg/api:** rename verifyRow request struct
- **pkg/api:** verified SQL get -> verified row renaming

### Features
- add verified sql get and fix tests
- **pkg/gw:** add immuerror server mux error handler


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


[Unreleased]: https://github.com/vchain-us/immudb/compare/v1.3.1...HEAD
[v1.3.1]: https://github.com/vchain-us/immudb/compare/v1.3.0...v1.3.1
[v1.3.0]: https://github.com/vchain-us/immudb/compare/v1.3.0-RC1...v1.3.0
[v1.3.0-RC1]: https://github.com/vchain-us/immudb/compare/v1.2.0...v1.3.0-RC1
[v1.2.0]: https://github.com/vchain-us/immudb/compare/v1.0.5...v1.2.0
[v1.0.5]: https://github.com/vchain-us/immudb/compare/v0.9.2-RC1...v1.0.5
[v0.9.2-RC1]: https://github.com/vchain-us/immudb/compare/v0.8.1...v0.9.2-RC1
