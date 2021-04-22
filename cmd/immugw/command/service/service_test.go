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

package service

import (
	"bytes"
	"github.com/codenotary/immudb/cmd/helper"
	"github.com/codenotary/immudb/cmd/immudb/command/service/servicetest"
	"io/ioutil"
	"testing"

	"github.com/codenotary/immudb/pkg/client/clienttest"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCommandLine_ServiceImmugwInstall(t *testing.T) {
	cmd := &cobra.Command{}
	tr := &clienttest.TerminalReaderMock{}
	tr.Responses = []string{"y"}
	cld := commandline{helper.Config{}, servicetest.NewSservicemock(), tr}

	cld.Service(cmd)
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"service", "install"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestCommandLine_ServiceImmugwUninstallAbortUnintall(t *testing.T) {
	cmd := &cobra.Command{}
	tr := &clienttest.TerminalReaderMock{}
	tr.Responses = []string{"n"}
	cld := commandline{helper.Config{}, servicetest.NewSservicemock(), tr}

	cld.Service(cmd)
	cmd.SetArgs([]string{"service", "uninstall"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestCommandLine_ServiceImmugwUninstallRemovingData(t *testing.T) {
	cmd := &cobra.Command{}
	tr := &clienttest.TerminalReaderMock{}
	tr.Responses = []string{"y", "y"}
	cld := commandline{helper.Config{}, servicetest.NewSservicemock(), tr}

	cld.Service(cmd)
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"service", "uninstall"})
	err := cmd.Execute()
	assert.Nil(t, err)
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(out), "uninstall")
}

func TestCommandLine_ServiceImmugwUninstallWithoutRemoveData(t *testing.T) {
	cmd := &cobra.Command{}
	tr := &clienttest.TerminalReaderMock{}
	tr.Responses = []string{"y", "n"}
	cld := commandline{helper.Config{}, servicetest.NewSservicemock(), tr}

	cld.Service(cmd)
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"service", "uninstall"})
	err := cmd.Execute()
	assert.Nil(t, err)
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(out), "uninstall")
}

func TestCommandLine_ServiceImmugwStop(t *testing.T) {
	cmd := &cobra.Command{}
	tr := &clienttest.TerminalReaderMock{}
	cld := commandline{helper.Config{}, servicetest.NewSservicemock(), tr}
	cld.Service(cmd)
	cmd.SetArgs([]string{"service", "stop"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestCommandLine_ServiceImmugwStart(t *testing.T) {
	cmd := &cobra.Command{}
	tr := &clienttest.TerminalReaderMock{}
	cld := commandline{helper.Config{}, servicetest.NewSservicemock(), tr}
	cld.Service(cmd)
	cmd.SetArgs([]string{"service", "start"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestCommandLine_ServiceImmugwDelayed(t *testing.T) {
	cmd := &cobra.Command{}
	tr := &clienttest.TerminalReaderMock{}
	cld := commandline{helper.Config{}, servicetest.NewSservicemock(), tr}
	cld.Service(cmd)
	cmd.SetArgs([]string{"service", "stop", "--time", "20"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestCommandLine_ServiceImmugwRestart(t *testing.T) {
	cmd := &cobra.Command{}
	tr := &clienttest.TerminalReaderMock{}
	cld := commandline{helper.Config{}, servicetest.NewSservicemock(), tr}
	cld.Service(cmd)
	cmd.SetArgs([]string{"service", "restart"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestCommandLine_ServiceImmugwStatus(t *testing.T) {
	cmd := &cobra.Command{}
	tr := &clienttest.TerminalReaderMock{}
	cld := commandline{helper.Config{}, servicetest.NewSservicemock(), tr}
	cld.Service(cmd)
	cmd.SetArgs([]string{"service", "status"})
	err := cmd.Execute()
	assert.Nil(t, err)
}
