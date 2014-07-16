/*
   Copyright Evan Hazlett

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

/*
   These are structs from the Docker package to prevent the dependency
   on the Docker library.  We just use the structs so there is no need
   to bring in the external dependencies (libdevmapper, btrfs, etc.) just
   to use those.

*/
package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/shipyard/shipyard-agent/utils"
)

type (
	Port        string
	PortBinding struct {
		HostIp   string
		HostPort string
	}
	PortSet map[Port]struct{}

	PortMap map[Port][]PortBinding

	InfoPort struct {
		PrivatePort int
		PublicPort  int
		Type        string
	}

	State struct {
		sync.RWMutex
		Running    bool
		Pid        int
		ExitCode   int
		StartedAt  time.Time
		FinishedAt time.Time
		Ghost      bool
	}

	APIContainer struct {
		Id     string
		Image  string
		Create int
		Status string
		Ports  []InfoPort
	}
	Container struct {
		Id              string
		Args            []string
		Name            string
		Config          ContainerConfig
		Created         time.Time
		Driver          string
		HostConfig      HostConfig
		HostnamePath    string
		HostsPath       string
		Image           string
		NetworkSettings NetworkSettings
		Path            string
		ResolvConfPath  string
		State           State
		Volumes         map[string]string
	}

	ContainerConfig struct {
		AttachStderr    bool
		AttachStdin     bool
		AttachStdout    bool
		Cmd             []string
		Entrypoint      []string
		CpuShares       int64
		Dns             string
		Domainname      string
		Env             []string
		ExposedPorts    map[Port]struct{}
		Hostname        string
		Image           string
		Memory          float64
		MemorySwap      float64
		NetworkDisabled bool
		OnBuild         []string
		OpenStdin       bool
		PortSpecs       []string
		StdinOnce       bool
		Tty             bool
		User            string
		Volumes         map[string]struct{}
		VolumesFrom     string
		WorkingDir      string
	}

	KeyValuePair struct {
		Key   string
		Value string
	}

	HostConfig struct {
		Binds           []string
		ContainerIDFile string
		LxcConf         []KeyValuePair
		Privileged      bool
		PortBindings    PortMap
		Links           []string
		PublishAllPorts bool
	}

	PortMapping map[string]string

	NetworkSettings struct {
		IPAddress   string
		IPPrefixLen int
		Gateway     string
		Bridge      string
		PortMapping map[string]PortMapping
		Ports       PortMap
	}
)

// Returns a new mux subrouter that acts as an adapter to support the Docker API
func NewDockerSubrouter(router *mux.Router) *mux.Router {
	rtr := router.PathPrefix("/{apiVersion:v1.*}").Subrouter()
	rtr.HandleFunc("/{.*}", dockerHandler).Methods("GET", "PUT", "POST", "DELETE")
	return rtr
}

// Docker: generic handler
func dockerHandler(w http.ResponseWriter, req *http.Request) {
	utils.ProxyLocalDockerRequest(w, req, dockerURL)
}
