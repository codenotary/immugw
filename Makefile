# Copyright 2019-2020 vChain, Inc. 											\
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

VERSION=0.7.1
TARGETS=linux/amd64 windows/amd64 darwin/amd64 linux/s390x
SERVICE_EXE=${SERVICE_NAME}-v${VERSION}-windows-amd64.exe

PWD = $(shell pwd)
GO ?= go
DOCKER ?= docker
PROTOC ?= protoc
STRIP = strip

V_COMMIT := $(shell git rev-parse HEAD)
#V_BUILT_BY := "$(shell echo "`git config user.name`<`git config user.email`>")"
V_BUILT_BY := $(shell git config user.email)
V_BUILT_AT := $(shell date +%s)
V_LDFLAGS_COMMON := -X "github.com/codenotary/immudb/cmd/version.Version=$(VERSION)" \
					-X "github.com/codenotary/immudb/cmd/version.Commit=$(V_COMMIT)" \
					-X "github.com/codenotary/immudb/cmd/version.BuiltBy=$(V_BUILT_BY)"\
					-X "github.com/codenotary/immudb/cmd/version.BuiltAt=$(V_BUILT_AT)"
V_LDFLAGS_STATIC := ${V_LDFLAGS_COMMON} \
				  -X github.com/codenotary/immudb/cmd/version.Static=static \
				  -extldflags "-static"

.PHONY: all
all: immugw
	@echo 'Build successful, now you can launch immugw.'

.PHONY: rebuild
rebuild: clean build/codegen all

.PHONY: immugw
immugw:
	$(GO) build -v -ldflags '$(V_LDFLAGS_COMMON)' ./cmd/immugw

.PHONY: immugw-static
immugw-static:
	CGO_ENABLED=0 $(GO) build -a -tags netgo -ldflags '$(V_LDFLAGS_STATIC) -extldflags "-static"' ./cmd/immugw

# To view coverage as HTML run: go tool cover -html=coverage.txt
.PHONY: test
test:
	$(GO) vet ./...
	$(GO) test -v -failfast $(go list ./... | grep -v test) --race -coverprofile=coverage.txt -covermode=atomic ./...


.PHONY: clean
clean:
	rm -f immugw

.PHONY: man
man:
	$(GO) run ./cmd/immugw mangen ./cmd/docs/man/immugw

########################## releases scripts ############################################################################
.PHONY: CHANGELOG.md
CHANGELOG.md:
	git-chglog -o CHANGELOG.md

.PHONY: CHANGELOG.md.next-tag
CHANGELOG.md.next-tag:
	git-chglog -o CHANGELOG.md --next-tag v${VERSION}

.PHONY: build/xgo
build/xgo:
	$(DOCKER) build \
			-f ./build/xgo/Dockerfile \
			-t ${SERVICE_NAME}-xgo \
			--pull=true \
			./build/xgo

.PHONY: build/makensis
build/makensis:
	$(DOCKER) build \
		-f ./build/makensis/Dockerfile \
		-t ${SERVICE_NAME}-makensis \
		./build/makensis

.PHONY: clean/dist
clean/dist:
	rm -Rf ./dist

.PHONY: dist
dist: clean/dist build/xgo
	mkdir -p dist
	$(GO) build -a -tags netgo -ldflags '${V_LDFLAGS_STATIC}' \
			-o ./dist/${SERVICE_NAME}-v${VERSION}-linux-amd64-static \
     		./cmd/${SERVICE_NAME}
	$(DOCKER) run --rm \
			-v "${PWD}/dist:/dist" \
			-v "${PWD}:/source:ro" \
			-e GO111MODULE=on \
			-e FLAG_LDFLAGS="-s ${V_LDFLAGS_COMMON}" \
			-e TARGETS="${TARGETS}" \
			-e PACK=cmd/${SERVICE_NAME} \
			-e OUT=${SERVICE_NAME}-v${VERSION} \
			${SERVICE_NAME}-xgo .
	mv ./dist/${SERVICE_NAME}-v${VERSION}-windows-4.0-amd64.exe ./dist/${SERVICE_EXE}
	mv ./dist/${SERVICE_NAME}-v${VERSION}-darwin-10.6-amd64 ./dist/${SERVICE_NAME}-v${VERSION}-darwin-amd64

.PHONY: dist/${SERVICE_EXE}
dist/${SERVICE_EXE}:
	echo ${SIGNCODE_PVK_PASSWORD} | $(DOCKER) run --rm -i \
		-v ${PWD}/dist:/dist \
		-v ${SIGNCODE_SPC}:/certs/f.spc:ro \
		-v ${SIGNCODE_PVK}:/certs/f.pvk:ro \
		mono:6.8.0 signcode \
		-spc /certs/f.spc -v /certs/f.pvk \
		-a sha1 -$ commercial \
		-n "CodeNotary ${SERVICE_NAME}" \
		-i https://codenotary.io/ \
		-t http://timestamp.comodoca.com -tr 10 \
		$@
	rm -Rf $@.bak

.PHONY: dist/sign
dist/sign: vendor ${SERVICE_NAME}
	for f in ./dist/*; do vcn sign -p $$f; printf "\n\n"; done

# SERVICE_NAME=immugw SIGNCODE_PVK_PASSWORD=<pvk password> SIGNCODE_PVK=<path to vchain.pvk> SIGNCODE_SPC=<path to vchain.spc> make dist/all
.PHONY: dist/all
dist/all: dist dist/${SERVICE_EXE}

.PHONY: dist/binary.md
dist/binary.md:
	@for f in ./dist/*; do \
		ff=$$(basename $$f); \
		shm_id=$$(sha256sum $$f | awk '{print $$1}'); \
		printf "[$$ff](https://github.com/vchain-us/immudb/releases/download/v${VERSION}/$$ff) | $$shm_id \n" ; \
	done
########################## releases scripts end ########################################################################
