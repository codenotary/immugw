# Copyright 2021 CodeNotary, Inc. All rights reserved. 											\
																			\
Licensed under the Apache License, Version 2.0 (the "License"); 			\
you may not use this file except in compliance with the License. 			\
You may obtain a copy of the License at 									\
																			\
	http://www.apache.org/licenses/LICENSE-2.0 								\
																			\
Unless required by applicable law or agreed to in writing, software 		\
distributed under the License is distributed on an "AS IS" BASIS, 			\
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.	\
See the License for the specific language governing permissions and 		\
limitations under the License.

export GO111MODULE=on

SHELL=/bin/bash -o pipefail

VERSION=1.1.0
SERVICES=immugw
TARGETS=linux/amd64 windows/amd64 darwin/amd64 linux/s390x linux/arm64 freebsd/amd64

PWD = $(shell pwd)
GO ?= go
GOPATH ?= $(shell go env GOPATH)
DOCKER ?= docker
PROTOC ?= protoc
STRIP = strip

V_COMMIT := $(shell git rev-parse HEAD)
V_BUILT_BY := $(shell git config user.email)
V_BUILT_AT := $(shell date +%s)
V_LDFLAGS_COMMON := -s -X "github.com/codenotary/immugw/cmd/version.Version=${VERSION}" \
					-X "github.com/codenotary/immugw/cmd/version.Commit=${V_COMMIT}" \
					-X "github.com/codenotary/immugw/cmd/version.BuiltBy=${V_BUILT_BY}"\
					-X "github.com/codenotary/immugw/cmd/version.BuiltAt=${V_BUILT_AT}"
V_LDFLAGS_STATIC := ${V_LDFLAGS_COMMON} \
				  -X github.com/codenotary/immugw/cmd/version.Static=static \
				  -extldflags "-static"

.PHONY: all
all: immugw
	@echo 'Build successful, now you can launch immugw.'

.PHONY: immugw
immugw:
	$(GO) build $(IMMUDB_BUILD_TAGS) -v -ldflags '$(V_LDFLAGS_COMMON)' ./cmd/immugw


.PHONY: immugw-static
immugw-static:
	CGO_ENABLED=0 $(GO) build $(IMMUDB_BUILD_TAGS) -a -ldflags '$(V_LDFLAGS_STATIC) -extldflags "-static"' ./cmd/immugw


.PHONY: test
test:
	$(GO) vet ./...
	$(GO) test -failfast --race ./...

.PHONY: clean
clean:
	rm -rf immugw
########################## releases scripts ############################################################################
.PHONY: CHANGELOG.md
CHANGELOG.md:
	git-chglog -o CHANGELOG.md

.PHONY: CHANGELOG.md.next-tag
CHANGELOG.md.next-tag:
	git-chglog -o CHANGELOG.md --next-tag v${VERSION}

.PHONY: clean/dist
clean/dist:
	rm -Rf ./dist

.PHONY: dist/binaries
dist/binaries:
		mkdir -p dist; \
		for service in ${SERVICES}; do \
    		for os_arch in ${TARGETS}; do \
    			goos=`echo $$os_arch|sed 's|/.*||'`; \
    			goarch=`echo $$os_arch|sed 's|^.*/||'`; \
    			printf "building $$goos $$goarch \n"; \
    		    GOOS=$$goos GOARCH=$$goarch $(GO) build -v -ldflags '${V_LDFLAGS_COMMON}' -o ./dist/$$service-v${VERSION}-$$goos-$$goarch ./cmd/$$service/$$service.go ; \
    		done; \
    		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -a -ldflags '${V_LDFLAGS_STATIC} -extldflags "-static"' -o ./dist/$$service-v${VERSION}-linux-amd64-static ./cmd/$$service/$$service.go ; \
    		mv ./dist/$$service-v${VERSION}-windows-amd64 ./dist/$$service-v${VERSION}-windows-amd64.exe; \
    	done

.PHONY: dist/sign
dist/sign:
	for f in ./dist/*; do vcn sign -p $$f; printf "\n\n"; done


.PHONY: dist/binary.md
dist/binary.md:
	@for f in ./dist/*; do \
		ff=$$(basename $$f); \
		shm_id=$$(sha256sum $$f | awk '{print $$1}'); \
		printf "[$$ff](https://github.com/codenotary/immugw/releases/download/v${VERSION}/$$ff) | $$shm_id \n" ; \
	done

########################## releases scripts end ########################################################################
