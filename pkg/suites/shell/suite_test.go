// Copyright (c) 2021 Doc.ai and/or its affiliates.
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

package shell_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
)

func TestShellFirstTry(t *testing.T) {
	t.Cleanup(func() { goleak.VerifyNone(t) })

	tempDir := t.TempDir()

	suite := shell.Suite{}
	suite.SetT(t)
	r := suite.Runner(tempDir)

	fileContent := "file content"
	fileName := "TestShellFirstTry.file"

	r.Run("echo " + fileContent + " >" + fileName)
	bytes, err := os.ReadFile(filepath.Clean(filepath.Join(tempDir, fileName)))
	require.NoError(t, err)
	require.Equal(t, fileContent+"\n", string(bytes))
}

func TestShellEventually(t *testing.T) {
	t.Cleanup(func() { goleak.VerifyNone(t) })

	tempDir := t.TempDir()

	suite := shell.Suite{}
	suite.SetT(t)
	r := suite.Runner(tempDir)

	fileName := "TestShellEventually.file"

	r.Run("X=1")
	r.Run(`if ! [[ $X == "1111" ]]; then
	echo $X >>` + fileName + `
	X=${X}1
	echo >&2 not enough ones
	false
fi`)
	bytes, err := os.ReadFile(filepath.Clean(filepath.Join(tempDir, fileName)))
	require.NoError(t, err)
	require.Equal(t, "1\n11\n111\n", string(bytes))
}
