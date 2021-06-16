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

package bash_test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/networkservicemesh/gotestmd/pkg/bash"
)

func TestBashProc(t *testing.T) {
	t.Cleanup(func() { goleak.VerifyNone(t) })

	runner, err := bash.New()
	require.NoError(t, err)
	defer runner.Close()

	stdout, stderr, exitCode, err := runner.Run("A=hello")
	require.NoError(t, err)
	require.Zero(t, exitCode)
	require.Empty(t, stdout)
	require.Empty(t, stderr)

	stdout, stderr, exitCode, err = runner.Run("B=world")
	require.NoError(t, err)
	require.Zero(t, exitCode)
	require.Empty(t, stdout)
	require.Empty(t, stderr)

	stdout, stderr, exitCode, err = runner.Run("echo $A $B")
	require.NoError(t, err)
	require.Zero(t, exitCode)
	require.Equal(t, "hello world", stdout)
	require.Empty(t, stderr)

	stdout, stderr, exitCode, err = runner.Run("abcdefg")
	require.NoError(t, err)
	require.NotZero(t, exitCode)
	require.Empty(t, stdout)
	require.Contains(t, stderr, "command not found")
}

func TestBashWriteFile(t *testing.T) {
	t.Cleanup(func() { goleak.VerifyNone(t) })

	runner, err := bash.New()
	require.NoError(t, err)
	defer runner.Close()

	envValue := "ns-1"

	stdout, stderr, exitCode, err := runner.Run("NAMESPACE=" + envValue)
	require.NoError(t, err)
	require.Zero(t, exitCode)
	require.Empty(t, stdout)
	require.Empty(t, stderr)

	stdout, stderr, exitCode, err = runner.Run(`cat > test <<EOF
$NAMESPACE
EOF`)
	require.NoError(t, err)
	require.Zero(t, exitCode)
	require.Empty(t, stdout)
	require.Empty(t, stderr)

	content, err := ioutil.ReadFile("test")
	require.NoError(t, err)
	require.Equal(t, envValue+"\n", string(content))
	_ = os.Remove("test")
}

func TestBashLongOperation(t *testing.T) {
	t.Cleanup(func() { goleak.VerifyNone(t) })

	runner, err := bash.New()
	require.NoError(t, err)
	defer runner.Close()

	stdout, stderr, exitCode, err := runner.Run("sleep 1s; echo hi")
	require.NoError(t, err)
	require.Zero(t, exitCode)
	require.Equal(t, "hi", stdout)
	require.Empty(t, stderr)
}

func TestBashMultilineOutput(t *testing.T) {
	t.Cleanup(func() { goleak.VerifyNone(t) })

	runner, err := bash.New()
	require.NoError(t, err)
	defer runner.Close()

	var text string
	for i := 0; i < 100; i++ {
		text += randomString(50) + "\n"
	}

	stdout, stderr, exitCode, err := runner.Run("echo -n $'" + text + "'")
	require.NoError(t, err)
	require.Zero(t, exitCode)
	require.Empty(t, stderr)

	// runner deleted the last '\n'
	require.Equal(t, len(text)-1, len(stdout))
	require.Equal(t, text[:len(text)-1], stdout)
}

func TestBashStderr(t *testing.T) {
	t.Cleanup(func() { goleak.VerifyNone(t) })

	runner, err := bash.New()
	require.NoError(t, err)
	defer runner.Close()

	stdout, stderr, exitCode, err := runner.Run(`echo out
echo err >&2`)
	require.NoError(t, err)
	require.Zero(t, exitCode)
	require.Equal(t, "out", stdout)
	require.Equal(t, "err", stderr)
}

func TestBashExitCode(t *testing.T) {
	t.Cleanup(func() { goleak.VerifyNone(t) })

	runner, err := bash.New()
	require.NoError(t, err)
	defer runner.Close()

	stdout, stderr, exitCode, err := runner.Run(`$(exit 42)`)
	require.NoError(t, err)
	require.Equal(t, 42, exitCode)
	require.Empty(t, stdout)
	require.Empty(t, stderr)
}

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		// #nosec
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
