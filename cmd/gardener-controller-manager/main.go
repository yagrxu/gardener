// Copyright 2018 The Gardener Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gardener/gardener/cmd/gardener-controller-manager/app"
)

func main() {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	// Setup signal handler if running inside a Kubernetes cluster
	stopCh := make(chan struct{})
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		c := make(chan os.Signal, 2)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			<-c
			close(stopCh)
			<-c
			os.Exit(1)
		}()
	}

	command := app.NewCommandStartGardenerControllerManager(stopCh)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
