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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/codenotary/immudb/pkg/api/schema"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/server"
	"github.com/codenotary/immudb/pkg/server/servertest"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func newTestClient(t *testing.T, dbs []string) Client {
	options := server.DefaultOptions().WithAuth(true).WithDir(t.TempDir())
	bs := servertest.NewBufconnServer(options)

	bs.Start()

	t.Cleanup(func() {
		bs.Stop()
		matches, _ := filepath.Glob("state-*")
		os.RemoveAll(options.Dir)
		for _, m := range matches {
			os.RemoveAll(m)
		}
	})

	client, err := bs.NewAuthenticatedClient(immuclient.
		DefaultOptions().
		WithDir(t.TempDir()),
	)
	require.NoError(t, err)

	t.Cleanup(func() { client.CloseSession(context.Background()) })

	var (
		ctx = context.TODO()
	)

	_, err = client.Login(ctx, []byte("immudb"), []byte("immudb"))
	require.NoError(t, err)

	for _, db := range dbs {
		// step 1: create test database
		err := client.CreateDatabase(ctx, &schema.DatabaseSettings{DatabaseName: db})
		require.NoError(t, err)
	}

	opts := immuclient.DefaultOptions().WithDialOptions([]grpc.DialOption{grpc.WithContextDialer(bs.Dialer), grpc.WithInsecure()}).WithAuth(false).WithDir(t.TempDir())
	cli := New(opts)

	return cli
}

func Test_client_add(t *testing.T) {
	cli := newTestClient(t, nil)

	// list of databases to add
	dbs := []string{"foodb", "bazdb"}

	// check if adding a new db works
	for _, db := range dbs {
		_, err := cli.Add(db)
		require.NoError(t, err)
	}
	require.Equal(t, len(cli.(*client).dbMap), len(dbs))

	// check if getting a db works
	for _, db := range dbs {
		c, err := cli.For(db)
		require.NoError(t, err)
		require.NotNil(t, c)
	}

	// check if getting an existing db works without adding it again to the dbmap
	for _, db := range dbs {
		_, err := cli.Add(db)
		require.NoError(t, err)
	}
	require.Equal(t, len(cli.(*client).dbMap), len(dbs))
}

func Test_client_concurrent_access(t *testing.T) {
	cli := newTestClient(t, nil)

	// list of databases to add
	dbs := []string{"foodb", "bazdb", "bardb", "casdb"}

	// check if adding a new db works
	for _, db := range dbs {
		_, err := cli.Add(db)
		require.NoError(t, err)
	}
	require.Equal(t, len(cli.(*client).dbMap), len(dbs))

	// check if concurrent access to the same db works
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			db := dbs[j%len(dbs)]

			c, err := cli.For(db)
			require.NoError(t, err)
			require.NotNil(t, c)
			require.Equal(t, db, c.GetOptions().Database)
		}(i)
	}
	wg.Wait()
}

func Test_client_concurrent_operation(t *testing.T) {
	// list of databases to add
	dbs := []string{"foodb", "bazdb", "bardb"}

	cli := newTestClient(t, dbs)

	// check if adding a new db works
	for _, db := range dbs {
		_, err := cli.Add(db)
		require.NoError(t, err)
	}
	require.Equal(t, len(cli.(*client).dbMap), len(dbs))

	// check if concurrent access to the same db works
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			db := dbs[j%len(dbs)]

			c, err := cli.For(db)
			require.NoError(t, err)
			require.NotNil(t, c)

			key := fmt.Sprintf("key%d", j)
			val := fmt.Sprintf("val%d", j)

			lr, err := c.Login(context.Background(), []byte("immudb"), []byte("immudb"))
			require.NoError(t, err)

			md := metadata.Pairs("authorization", lr.Token)
			testUserContext := metadata.NewOutgoingContext(context.Background(), md)

			dbResp, err := c.UseDatabase(testUserContext, &schema.Database{DatabaseName: db})
			md = metadata.Pairs("authorization", dbResp.Token)
			testUserContext = metadata.NewOutgoingContext(context.Background(), md)

			// verify if the correct db is selected and being used
			require.Equal(t, db, c.GetOptions().Database)
			require.Equal(t, db, c.GetOptions().CurrentDatabase)

			// set key
			_, err = c.Set(testUserContext, []byte(key), []byte(val))
			require.NoError(t, err)

			// get key
			resp, err := c.Get(testUserContext, []byte(key))
			require.NoError(t, err)
			require.Equal(t, resp.Value, []byte(val))
		}(i)
	}
	wg.Wait()

}

func Test_client_mocks(t *testing.T) {
	c := NewMockClient(nil, immuclient.DefaultOptions())
	require.Nil(t, c.(*client).dbMap["defaultdb"])

	c = NewMockClientWithDb(nil, immuclient.DefaultOptions(), "bazdb")
	require.Nil(t, c.(*client).dbMap["bazdb"])
}
