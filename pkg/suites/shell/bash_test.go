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

package shell_test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
)

func TestShellProc(t *testing.T) {
	var bash shell.Bash
	defer bash.Close()

	stdout, stderr, success, err := bash.Run("A=hello")
	require.NoError(t, err)
	require.True(t, success)
	require.Empty(t, stdout)
	require.Empty(t, stderr)

	stdout, stderr, success, err = bash.Run("B=world")
	require.NoError(t, err)
	require.True(t, success)
	require.Empty(t, stdout)
	require.Empty(t, stderr)

	stdout, stderr, success, err = bash.Run("echo $A $B")
	require.NoError(t, err)
	require.True(t, success)
	require.Equal(t, "hello world", stdout)
	require.Empty(t, stderr)

	stdout, stderr, success, err = bash.Run("abcdefg")
	require.NoError(t, err)
	require.False(t, success)
	require.Empty(t, stdout)
	require.Contains(t, stderr, "command not found")
}
func TestShellWriteFile(t *testing.T) {
	var bash shell.Bash
	defer bash.Close()

	envValue := "ns-1"

	stdout, stderr, success, err := bash.Run("NAMESPACE=" + envValue)
	require.NoError(t, err)
	require.True(t, success)
	require.Empty(t, stdout)
	require.Empty(t, stderr)

	stdout, stderr, success, err = bash.Run(`cat > test <<EOF
$NAMESPACE
EOF`)
	require.NoError(t, err)
	require.True(t, success)
	require.Empty(t, stdout)
	require.Empty(t, stderr)

	content, err := ioutil.ReadFile("test")
	require.NoError(t, err)
	require.Equal(t, envValue+"\n", string(content))
	_ = os.Remove("test")
}

func TestShellLongOperation(t *testing.T) {
	var bash shell.Bash
	defer bash.Close()

	stdout, stderr, success, err := bash.Run("sleep 1s; echo hi")
	require.NoError(t, err)
	require.True(t, success)
	require.Equal(t, "hi", stdout)
	require.Empty(t, stderr)
}

func TestShellMultilineOutput(t *testing.T) {
	var bash shell.Bash
	defer bash.Close()

	// Generate 100 text lines. Each line has 50 characters + '\n', so 5100
	var text string
	for i := 0; i < 100; i++ {
		text += randomString(50) + "\n"
	}

	stdout, stderr, success, err := bash.Run("echo -n $'" + text + "'")
	require.NoError(t, err)
	require.True(t, success)
	require.Empty(t, stderr)

	// Bash deleted the last '\n'
	require.Equal(t, len(text)-1, len(stdout))
	require.Equal(t, text[:len(text)-1], stdout)
}

func TestShellStderr(t *testing.T) {
	var bash shell.Bash
	defer bash.Close()

	stdout, stderr, success, err := bash.Run(`echo out
echo err >&2`)
	require.NoError(t, err)
	require.True(t, success)
	require.Equal(t, "out", stdout)
	require.Equal(t, "err", stderr)
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
