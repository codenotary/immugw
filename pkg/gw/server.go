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
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/codenotary/immudb/pkg/client/auditor"
	"github.com/codenotary/immudb/pkg/client/cache"
	"github.com/codenotary/immudb/pkg/client/state"
	"github.com/codenotary/immugw/pkg/api"
	immugwclient "github.com/codenotary/immugw/pkg/client"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/immuos"
	"github.com/codenotary/immudb/pkg/server"
	"github.com/codenotary/immugw/pkg/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
)

var startedAt time.Time

// Start starts the immudb gateway server
func (s *ImmuGwServer) Start() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := immugwclient.New(&s.CliOptions)

	ic, err := client.Add("defaultdb") // TODO: fix this and make this dynamic
	if err != nil {
		s.Logger.Errorf("unable to instantiate client: %s", err)
		return err
	}

	ic.WithTokenService(s.Options.TokenService)

	mux := runtime.NewServeMux(runtime.WithProtoErrorHandler(api.DefaultGWErrorHandler))

	handler := cors.Default().Handler(mux)

	rt := DefaultRuntime()
	json := json.DefaultJSON()

	sh := NewSetHandler(mux, client, rt, json)
	ssh := NewVerifiedSetHandler(mux, client, rt, json)
	sgh := NewVerifiedGetHandler(mux, client, rt, json)
	hh := NewHistoryHandler(mux, client, rt, json)
	sr := NewSafeReferenceHandler(mux, client, rt, json)
	sza := NewVerifiedZaddHandler(mux, client, rt, json)
	udb := NewUseDatabaseHandler(mux, client, rt, json)
	tx := NewVerifiedTxByIdHandler(mux, client, rt, json)
	vsql := NewVerifiedSQLGetHandler(mux, client, rt, json)

	mux.Handle(http.MethodPost, api.Pattern_ImmuService_Set_0, sh.Set)
	mux.Handle(http.MethodPost, api.Pattern_ImmuService_VerifiedSet_0(), ssh.VerifiedSet)
	mux.Handle(http.MethodPost, api.Pattern_ImmuService_VerifiedGet_0(), sgh.VerifiedGet)
	mux.Handle(http.MethodPost, api.Pattern_ImmuService_History_0, hh.History)
	mux.Handle(http.MethodPost, api.Pattern_ImmuService_VerifiedSetReference_0(), sr.SafeReference)
	mux.Handle(http.MethodPost, api.Pattern_ImmuService_VerifiedZAdd_0(), sza.VerifiedZadd)
	mux.Handle(http.MethodGet, schema.Pattern_ImmuService_UseDatabase_0(), udb.UseDatabase)
	mux.Handle(http.MethodGet, api.Pattern_ImmuService_VerifiedTxById_0(), tx.VerifiedTxById)
	mux.Handle(http.MethodPost, api.Pattern_ImmuService_VerifiableSQLGet_0(), vsql.VerifiedSQLGetHandler)

	err = RegisterImmuServiceHandlerClient(ctx, mux, client, ic.GetServiceClient())
	if err != nil {
		s.Logger.Errorf("unable to register client handlers: %s", err)
		return err
	}

	s.installShutdownHandler()
	s.Logger.Infof("starting immugw: %v", s.Options)
	if s.Options.Pidfile != "" {
		if s.Pid, err = server.NewPid(s.Options.Pidfile, immuos.NewStandardOS()); err != nil {
			s.Logger.Errorf("failed to write pidfile: %s", err)
			return err
		}
	}

	if s.Options.Audit {
		defaultAuditor, err := auditor.DefaultAuditor(
			s.Options.AuditInterval,
			fmt.Sprintf("%s:%d", s.Options.ImmudbAddress, s.Options.ImmudbPort),
			s.CliOptions.DialOptions,
			s.Options.AuditUsername,
			s.Options.AuditPassword,
			nil,
			nil,
			auditor.AuditNotificationConfig{},
			ic.GetServiceClient(),
			state.NewUUIDProvider(ic.GetServiceClient()),
			cache.NewHistoryFileCache(filepath.Join(s.CliOptions.Dir, "auditor")),
			s.MetricServer.mc.UpdateAuditResult,
			s.Logger,
			nil)
		if err != nil {
			s.Logger.Errorf("unable to create auditor: %s", err)
			return err
		}
		go defaultAuditor.Run(s.Options.AuditInterval, false, ctx.Done(), s.auditorDone)
		defer func() { <-s.auditorDone }()
	}

	go func() {
		if err = http.ListenAndServe(s.Options.Address+":"+strconv.Itoa(s.Options.Port), handler); err != nil && err != http.ErrServerClosed {
			s.Logger.Errorf("unable to launch immugw: %+s", err)
		}
	}()

	metricsServer := s.MetricServer.StartMetrics()
	defer func() {
		if err = metricsServer.Close(); err != nil {
			s.Logger.Errorf("failed to shutdown metric server: %s", err)
		}
	}()
	startedAt = time.Now()
	<-s.quit
	if s.Options.Audit {
		cancel()
	}
	return err
}

// Stop stops the immudb gateway server
func (s *ImmuGwServer) Stop() error {
	s.Logger.Infof("stopping immugw: %v", s.Options)
	defer func() { s.quit <- struct{}{} }()
	return nil
}

func (s *ImmuGwServer) installShutdownHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer func() {
			s.quit <- struct{}{}
		}()
		<-c
		s.Logger.Debugf("Caught SIGTERM")
		if err := s.Stop(); err != nil {
			s.Logger.Errorf("Shutdown error: %v", err)
		}
		s.Logger.Infof("Shutdown completed")
	}()
}
