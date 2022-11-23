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

package client

import (
	"fmt"
	"sync"

	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// api errors
var (
	ErrDatabaseNotFound = status.Error(codes.NotFound, "database is not initialised")
)

// New returns a new Client using the Options to connect to immudb.
// The returned client reads *and* writes directly to the server
// It understands how to work with with multiple databases.
func New(options *immuclient.Options) Client {
	return newClient(options)
}

func newClient(opts *immuclient.Options) *client {
	return &client{
		opts:   opts,
		dbList: make(map[string]immuclient.ImmuClient),
	}
}

type client struct {
	mu     sync.RWMutex
	opts   *immuclient.Options
	dbList map[string]immuclient.ImmuClient
}

// Add adds a client connection for database db to the immudb server
func (c *client) Add(db string) (immuclient.ImmuClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	dir := fmt.Sprintf("state-%s", db)
	opts := *c.opts
	opts.WithDir(dir)

	cli, err := immuclient.NewImmuClient(&opts)
	if err != nil {
		return nil, err
	}
	c.dbList[db] = cli
	return cli, nil
}

// For returns the client for database db to the immudb server
func (c *client) For(db string) (immuclient.ImmuClient, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.dbList[db]
	if !ok {
		return nil, ErrDatabaseNotFound
	}
	return v, nil
}

// NewMockClient returns a mock Client for defaultdb to the immudb server
func NewMockClient(cli immuclient.ImmuClient, opts *immuclient.Options) Client {
	return &client{
		opts: opts,
		dbList: map[string]immuclient.ImmuClient{
			"defaultdb": cli,
		},
	}
}

// NewMockClientWithDb returns a mock Client for database db to the immudb server
func NewMockClientWithDb(cli immuclient.ImmuClient, opts *immuclient.Options, db string) Client {
	return &client{
		opts: opts,
		dbList: map[string]immuclient.ImmuClient{
			db: cli,
		},
	}
}
