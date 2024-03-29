/*
Copyright 2019-2020 vChain, Inc.

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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/codenotary/immudb/embedded/logger"
	"github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/client/homedir"
	"github.com/codenotary/immudb/pkg/client/tokenservice"
	"github.com/codenotary/immudb/pkg/server"
	"github.com/codenotary/immudb/pkg/server/servertest"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestImmuGwServer_Start(t *testing.T) {
	options := server.DefaultOptions().WithAuth(true)
	bs := servertest.NewBufconnServer(options)

	bs.Start()
	defer bs.Stop()

	defer func() {
		matches, _ := filepath.Glob("state-*")
		os.RemoveAll(options.Dir)
		for _, m := range matches {
			os.RemoveAll(m)
		}
	}()

	cliOpts := client.DefaultOptions().
		WithDir(options.Dir).
		WithPort(0).
		WithAddress("").
		WithHealthCheckRetries(1).
		WithMTLsOptions(client.MTLsOptions{}).
		WithMaxRecvMsgSize(4 * 1024 * 1024).
		WithAuth(true).
		WithConfig("").
		WithDialOptions([]grpc.DialOption{
			grpc.WithContextDialer(bs.Dialer), grpc.WithInsecure(),
		})

	l := logger.NewSimpleLogger("test", os.Stdout)

	gwOpts := Options{}.WithTokenService(tokenservice.NewFileTokenService().WithHds(homedir.NewHomedirService()))

	gw := ImmuGwServer{
		Options:      gwOpts,
		CliOptions:   *cliOpts,
		Logger:       l,
		quit:         make(chan struct{}, 1),
		MetricServer: newMetricsServer(DefaultOptions().MetricsBind(), l, func() float64 { return time.Since(startedAt).Hours() }),
	}
	gw.quit <- struct{}{}
	err := gw.Start()
	assert.Nil(t, err)
}

func TestImmuGwServer_StartWithAuditor(t *testing.T) {
	options := server.DefaultOptions().WithAuth(true)
	bs := servertest.NewBufconnServer(options)

	bs.Start()
	defer bs.Stop()

	defer func() {
		matches, _ := filepath.Glob("state-*")
		os.RemoveAll(options.Dir)
		for _, m := range matches {
			os.RemoveAll(m)
		}
	}()

	cliOpts := client.DefaultOptions().
		WithDir(options.Dir).
		WithPort(0).
		WithAddress("").
		WithHealthCheckRetries(1).
		WithMTLsOptions(client.MTLsOptions{}).
		WithMaxRecvMsgSize(4 * 1024 * 1024).
		WithAuth(true).
		WithConfig("").
		WithDialOptions([]grpc.DialOption{
			grpc.WithContextDialer(bs.Dialer), grpc.WithInsecure(),
		})

	l := logger.NewSimpleLogger("test", os.Stdout)

	gwOpts := Options{}.WithAudit(true).
		WithAuditInterval(5 * time.Millisecond).
		WithTokenService(tokenservice.NewFileTokenService().WithHds(homedir.NewHomedirService()))

	gw := ImmuGwServer{
		Options:      gwOpts,
		CliOptions:   *cliOpts,
		Logger:       l,
		quit:         make(chan struct{}, 1),
		auditorDone:  make(chan struct{}, 1),
		MetricServer: newMetricsServer(DefaultOptions().MetricsBind(), l, func() float64 { return time.Since(startedAt).Hours() }),
	}
	gw.quit <- struct{}{}
	gw.auditorDone <- struct{}{}

	err := gw.Start()
	assert.Nil(t, err)
}
