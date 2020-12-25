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

const emptyTest = `func (s *Suite) Test() {}`

const testTemplate = `
func (s *Suite) Test{{ NAME }}() {
	dir := filepath.Join(os.Getenv("GOPATH"), "src", "{{ DIR }}")
    r := s.Runner(dir)    
	{{ CLEANUP }}
	{{ RUN }}
}
`

// TestTemplate is a template for a test for a suite
type TestTemplate struct {
	Dir     string
	Name    string
	Cleanup Body
	Run     Body
}

// String returns string as a test for the suite
func (t *TestTemplate) String() string {
	template := testTemplate
	if len(t.Cleanup)+len(t.Run) == 0 {
		template = emptyTest
	}
	result := strings.ReplaceAll(template, "{{ NAME }}", t.Name)
	result = strings.ReplaceAll(result, "{{ DIR }}", t.Dir)

	cleanup := t.Cleanup.String()
	if len(cleanup) > 0 {
		cleanup = fmt.Sprintf(`	s.T().Cleanup(func() {
		%v
	})`, cleanup)
	}

	result = strings.ReplaceAll(result, "{{ CLEANUP }}", cleanup)
	result = strings.ReplaceAll(result, "{{ RUN }}", t.Run.String())

	return result
}
