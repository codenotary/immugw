/*
Copyright 2021 CodeNotary, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gw

import (
	"errors"
	"strings"

	"github.com/codenotary/immudb/pkg/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// sdk errors
var (
	ErrKeyNotFound   = status.Error(codes.Unknown, "key not found")
	ErrCorruptedData = status.Error(codes.Aborted, "data is corrupted") // codes.Aborted is translated in StatusConflict 409 http error
)

// wrap server errors which are not constants in immudb
var (
	ErrIllegalArgument   = errors.New("illegal arguments: empty key")
	ErrKeyNotFoundTBTree = errors.New("tbtree: key not found")
)

var (
	StatusErrKeyNotFound   = status.Error(codes.NotFound, "")
	StatusDatabaseNotFound = status.Error(codes.NotFound, "")
)

func mapSdkError(err error) error {
	switch {
	case errors.Is(err, ErrKeyNotFound):
		return StatusErrKeyNotFound
	case strings.HasPrefix(err.Error(), "data is corrupted"):
		return ErrCorruptedData
	case strings.HasSuffix(err.Error(), ErrIllegalArgument.Error()):
		return server.ErrIllegalArguments
	case strings.HasSuffix(err.Error(), ErrKeyNotFoundTBTree.Error()):
		return StatusErrKeyNotFound
	}
	return err
}
