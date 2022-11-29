/*
Copyright 2022 CodeNotary, Inc. All rights reserved.

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

var (
	// ErrDatabaseNotFound is returned when a database is not found
	ErrDatabaseNotFound = status.Error(codes.NotFound, "database is not initialised")
)

// New returns a new Client using the Options to connect to immudb.
// The returned client reads *and* writes directly to the server
// It understands how to work with with multiple databases.
func New(options *immuclient.Options) Client {
	return newClient(options)
}

// newClient returns a new Client for defaultdb to the immudb server
func newClient(opts *immuclient.Options) *client {
	return &client{
		opts:  opts,
		dbMap: make(map[string]immuclient.ImmuClient),
	}
}

// client implementa Client interface
type client struct {
	mu    sync.RWMutex
	opts  *immuclient.Options
	dbMap map[string]immuclient.ImmuClient
}

// Add adds a new database to the client
func (c *client) Add(db string) (immuclient.ImmuClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// check if db already exists
	if cli, ok := c.dbMap[db]; ok {
		return cli, nil
	}

	// create state dir for db
	dir := fmt.Sprintf("state-%s", db)
	opts := *c.opts
	opts.WithDir(dir).WithDatabase(db)

	// create new client
	cli, err := immuclient.NewImmuClient(&opts)
	if err != nil {
		return nil, err
	}

	// add client to map
	c.dbMap[db] = cli
	return cli, nil
}

// For returns the client connection for database db
func (c *client) For(db string) (immuclient.ImmuClient, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// check if db exists
	v, ok := c.dbMap[db]
	if !ok {
		return nil, ErrDatabaseNotFound
	}
	return v, nil
}

// NewMockClient returns a mock Client for defaultdb to the immudb server
func NewMockClient(cli immuclient.ImmuClient, opts *immuclient.Options) Client {
	return &client{
		opts: opts,
		dbMap: map[string]immuclient.ImmuClient{
			"defaultdb": cli,
		},
	}
}

// NewMockClientWithDb returns a mock Client for database db to the immudb server
func NewMockClientWithDb(cli immuclient.ImmuClient, opts *immuclient.Options, db string) Client {
	return &client{
		opts: opts,
		dbMap: map[string]immuclient.ImmuClient{
			db: cli,
		},
	}
}
