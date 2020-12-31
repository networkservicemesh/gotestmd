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

package templates

import (
	"fmt"
	"strings"
)

const suiteTemplate = `
package {{ NAME }}

import(
	"os"
	"path/filepath"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
	{{ IMPORTS }}
)

type Suite struct {
	shell.Suite
	{{ FIELDS }}
}

func (s *Suite) SetupSuite() {
	{{ SETUP }}
	dir := filepath.Join(os.Getenv("GOPATH"), "{{ DIR }}")
	r := s.Runner(dir)
	{{ CLEANUP }}
	{{ RUN }}
}
`

// Body represents a body of the method
type Body []string

// String returns the body as part of the method
func (b Body) String() string {
	var sb strings.Builder

	if len(b) == 0 {
		return ""
	}

	for _, block := range b {
		sb.WriteString("r.Run(")
		var lines = strings.Split(block, "\n")
		for i, line := range lines {
			sb.WriteString("`")
			sb.WriteString(line)
			sb.WriteString("`")
			if i+1 < len(lines) {
				sb.WriteString("+\"\\n\"+")
			}
		}
		sb.WriteString(")\n")
	}

	return sb.String()
}

// SuiteTemplate represents a template for generating a testify suite.Suite
type SuiteTemplate struct {
	Dir      string
	Location string
	Dependency
	Cleanup  Body
	Run      Body
	Tests    []*TestTemplate
	Deps     Dependencies
	TestDeps Dependencies
}

// String returns a string that contains generated testify.Suite
func (s *SuiteTemplate) String() string {
	result := strings.ReplaceAll(suiteTemplate, "{{ NAME }}", s.Name())
	cleanup := s.Cleanup.String()
	if len(cleanup) > 0 {
		cleanup = fmt.Sprintf(`	s.T().Cleanup(func() {
		%v
	})`, cleanup)
	}
	result = strings.ReplaceAll(result, "{{ CLEANUP }}", cleanup)
	result = strings.ReplaceAll(result, "{{ RUN }}", s.Run.String())
	result = strings.ReplaceAll(result, "{{ IMPORTS }}", s.Deps.String())
	result = strings.ReplaceAll(result, "{{ FIELDS }}", s.Deps.FieldsString())
	result = strings.ReplaceAll(result, "{{ SETUP }}", s.Deps.SetupString())
	result = strings.ReplaceAll(result, "{{ DIR }}", s.Dir)

	var sb strings.Builder
	_, _ = sb.WriteString(result)

	if len(s.Tests) == 0 {
		s.Tests = append(s.Tests, new(TestTemplate))
	}

	for _, test := range s.Tests {
		_, _ = sb.WriteString(test.String())
	}
	return spaceRegex.ReplaceAllString(strings.TrimSpace(sb.String()), "\n")
}
