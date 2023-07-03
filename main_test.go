// Copyright (c) 2020-2023 Doc.ai and/or its affiliates.
//
// Copyright (c) 2023 Cisco and/or its affiliates.
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

package main_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/networkservicemesh/gotestmd/pkg/bash"
)

func TestExamples(t *testing.T) {
	t.Cleanup(func() {
		_ = os.RemoveAll("test-examples")
	})
	runner, err := bash.New()
	require.NoError(t, err)
	defer runner.Close()
	_, _, exitCode, err := runner.Run("go install ./...")
	require.NoError(t, err)
	require.Zero(t, exitCode)

	_, _, exitCode, err = runner.Run("gotestmd examples/ test-examples/")
	require.NoError(t, err)
	require.Zero(t, exitCode)

	_, _, exitCode, err = runner.Run(`cat > test-examples/entry_point_test.go <<EOF
package suites

import (
	"testing"

	"github.com/networkservicemesh/gotestmd/test-examples/helloworld"
	"github.com/networkservicemesh/gotestmd/test-examples/producer/consumer2"
	"github.com/networkservicemesh/gotestmd/test-examples/producer/consumer3"
	"github.com/networkservicemesh/gotestmd/test-examples/tree"
	"github.com/stretchr/testify/suite"
)

func TestEntryPoint(t *testing.T) {
	suite.Run(t, new(helloworld.Suite))
	suite.Run(t, new(tree.Suite))
	suite.Run(t, new(consumer2.Suite))
	suite.Run(t, new(consumer3.Suite))
}
EOF
`)
	require.NoError(t, err)
	require.Zero(t, exitCode)

	_, _, exitCode, err = runner.Run("go test ./test-examples/... ")
	require.NoError(t, err)
	require.Zero(t, exitCode)
}

func TestBashSuite(t *testing.T) {
	t.Cleanup(func() {
		_ = os.RemoveAll("test-bash-examples")
	})
	runner, err := bash.New()
	require.NoError(t, err)
	defer runner.Close()
	_, _, exitCode, err := runner.Run("go install ./...")
	require.NoError(t, err)
	require.Zero(t, exitCode)

	_, _, exitCode, err = runner.Run("gotestmd examples/ test-bash-examples/ --bash --match=tree")
	require.NoError(t, err)
	require.Zero(t, exitCode)

	_, _, exitCode, err = runner.Run("./test-bash-examples/tree/suite.gen.sh setup")
	require.NoError(t, err)
	require.Zero(t, exitCode)

	_, _, exitCode, err = runner.Run("./test-bash-examples/tree/suite.gen.sh cleanup")
	require.NoError(t, err)
	require.Zero(t, exitCode)
}

func TestBashTest(t *testing.T) {
	t.Cleanup(func() {
		_ = os.RemoveAll("test-bash-examples")
	})
	runner, err := bash.New()
	require.NoError(t, err)
	defer runner.Close()
	_, _, exitCode, err := runner.Run("go install ./...")
	require.NoError(t, err)
	require.Zero(t, exitCode)

	_, _, exitCode, err = runner.Run("gotestmd examples/ test-bash-examples/ --bash --match=LeafA")
	require.NoError(t, err)
	require.Zero(t, exitCode)

	_, _, exitCode, err = runner.Run("./test-bash-examples/tree/suite.gen.sh setup")
	require.NoError(t, err)
	require.Zero(t, exitCode)

	_, _, exitCode, err = runner.Run("./test-bash-examples/tree/suite.gen.sh testLeafA")
	require.NoError(t, err)
	require.Zero(t, exitCode)

	_, _, exitCode, err = runner.Run("./test-bash-examples/tree/suite.gen.sh cleanup")
	require.NoError(t, err)
	require.Zero(t, exitCode)
}

func TestBashNoMatchesFound(t *testing.T) {
	t.Cleanup(func() {
		_ = os.RemoveAll("test-bash-examples")
	})
	runner, err := bash.New()
	require.NoError(t, err)
	defer runner.Close()
	_, _, exitCode, err := runner.Run("go install ./...")
	require.NoError(t, err)
	require.Zero(t, exitCode)

	_, stderr, exitCode, err := runner.Run("gotestmd examples/ test-bash-examples/ --bash --match=UnmatchablePattern")
	require.Contains(t, stderr, "No matches found for pattern: UnmatchablePattern")
	require.NoError(t, err)
	require.NotZero(t, exitCode)
}
