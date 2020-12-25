// Copyright (c) 2020 Doc.ai and/or its affiliates.
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
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
)

func TestShellProc(t *testing.T) {
	var bash shell.Bash
	defer bash.Close()

	_, err := bash.Run("A=hello")
	require.NoError(t, err)

	_, err = bash.Run("B=world")
	require.NoError(t, err)

	out, err := bash.Run("echo $A $B")
	require.NoError(t, err)
	require.Equal(t, "hello world", out)

	_, err = bash.Run("abcdefg")
	require.Error(t, err)
}
func TestShellWriteFile(t *testing.T) {
	var bash shell.Bash
	defer bash.Close()

	_, err := bash.Run("NAMESPACE=ns-1")
	require.NoError(t, err)
	_, err = bash.Run(`cat > test <<EOF
$NAMESPACE
EOF`)
	require.NoError(t, err)
	content, err := ioutil.ReadFile("test")
	require.NoError(t, err)
	require.Equal(t, "ns-1\n", string(content))
	_ = os.Remove("test")
}

func TestShellLongOperation(t *testing.T) {
	var bash shell.Bash
	defer bash.Close()

	out, err := bash.Run("sleep 1s; echo hi")
	require.NoError(t, err)
	require.Equal(t, "hi", out)
}
