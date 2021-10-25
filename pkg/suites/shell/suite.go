// Copyright (c) 2020-2021 Doc.ai and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package shell provides shell helpers and shell based suites
package shell

import (
	"flag"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/gotestmd/pkg/bash"
)

var timeoutFlag = flag.Duration("timeout", time.Minute, "timeout for command execution. Usage: set timeout in duratiom format via shell.timeout flag")
var once sync.Once

// Suite is testify suite that provides a shell helper functions for each test.
type Suite struct {
	suite.Suite
}

// Runner creates runner and sets the passed dir and envs
func (s *Suite) Runner(env ...string) *Runner {
	result := &Runner{
		t: s.T(),
	}
	b, err := bash.New(bash.WithEnv(env))
	if err != nil {
		s.FailNowf("can't initialize bash", "%v", err)
	}
	result.bash = b

	s.T().Cleanup(func() {
		result.bash.Close()
	})
	result.logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			DisableQuote: true,
		},
	}
	once.Do(func() {
		flag.Parse()
	})
	return result
}

// Runner is shell runner.
type Runner struct {
	t      *testing.T
	logger *logrus.Logger
	bash   *bash.Bash
}

// Run runs cmd, logs stdin, stdout, stderr
// Tries to run cmd several times, until it succeeds or timeout passes.
//
// Fails the test if the command can't be run successfully.
func (r *Runner) Run(cmd string) {
	timeoutCh := time.After(*timeoutFlag)
	for {
		r.logger.WithField(r.t.Name(), "stdin").Info(cmd)
		stdout, stderr, exitCode, err := r.bash.Run(cmd)
		if err != nil {
			r.logger.Fatalf("can't run command: %v", err)
			r.t.FailNow()
		}
		if stdout != "" {
			r.logger.WithField(r.t.Name(), "stdout").Info(stdout)
		}
		if stderr != "" {
			r.logger.WithField(r.t.Name(), "stderr").Info(stderr)
		}
		if exitCode == 0 {
			return
		}
		r.logger.WithField(r.t.Name(), "exitCode").Info(exitCode)
		select {
		case <-timeoutCh:
			r.logger.WithField("cmd", cmd).Error("command didn't succeed until timeout")
			require.Equal(r.t, 0, exitCode)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
