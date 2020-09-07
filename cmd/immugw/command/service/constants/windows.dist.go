// +build windows

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

package constants

import "fmt"

const ExecPath = ""
const ConfigPath = ""
const OSUser = ""
const OSGroup = ""

var StartUpConfig = ""

// UsageDet details on config and log file on specific os
var UsageDet = fmt.Sprintf(`Config and log files are present in C:\ProgramData\immugw folder`)

// UsageExamples usage examples for linux
var UsageExamples = fmt.Sprintf(`Install the immutable database
immugw.exe  service install    -  Initializes and runs daemon
immugw.exe  service stop       -  Stops the daemon
immugw.exe  service start      -  Starts initialized daemon
immugw.exe  service restart    -  Restarts daemon
immugw.exe  service uninstall  -  Removes daemon and its setup
Uninstall immugw after 20 second
immugw.exe  service uninstall --time 20 immugw.exe`)
