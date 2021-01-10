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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Suite is testify suite that provides a shell helper functions for each test.
// For each test generates a unique folder.
// Shell for each test located in the unique test folder.
type Suite struct {
	suite.Suite
}

// Runner creates runner and sets a passed dir and envs
func (s *Suite) Runner(dir string, env ...string) *Runner {
	result := &Runner{
		t: s.T(),
	}
	result.bash.Dir = filepath.Join(findRoot(), dir)
	result.bash.Env = env
	s.T().Cleanup(func() {
		result.bash.Close()
	})
	return result
}

func findRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err.Error())
	}
	currDir := wd
	for len(currDir) > 0 {
		if err != nil {
			logrus.Fatal(err.Error())
		}
		p := filepath.Clean(filepath.Join(currDir, "go.mod"))
		if _, err := os.Open(p); err == nil {
			return currDir
		}
		currDir = filepath.Dir(currDir)
	}
	return ""
}

// Runner is shell runner.
type Runner struct {
	t    *testing.T
	bash Bash
}

// Run runs cmd logs stdout, stderror, stdin
// Tries to run cmd on fail during timeout.
// Test could fail on the error or achieved cmd timeout.
func (r *Runner) Run(cmd string) {
	timeoutCh := time.After(time.Minute)
	var err error
	var out string
	for {
		logrus.WithField(r.t.Name(), "stdin").Info(cmd)
		out, err = r.bash.Run(cmd)
		if out != "" {
			logrus.WithField(r.t.Name(), "stdout").Info(out)
		}
		if err == nil {
			return
		}
		logrus.WithField(r.t.Name(), "stderr").Info(err.Error())
		select {
		case <-timeoutCh:
			require.NoError(r.t, err)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
