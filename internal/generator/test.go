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

package generator

import (
	"fmt"
	"strings"
	"text/template"
)

const emptyTest = `func (s *Suite) Test() {}`

const testTemplate = `
func (s *Suite) Test{{ .Name }}() {
	r := s.Runner()
	{{ .Cleanup }}
	{{ .Run }}
}
`

// Test is a template for a test for a suite
type Test struct {
	Dir     string
	Name    string
	Cleanup Body
	Run     Body
}

// String returns string as a test for the suite
func (t *Test) String() string {
	source := testTemplate
	if len(t.Cleanup)+len(t.Run) == 0 {
		source = emptyTest
	}

	tmpl, err := template.New("test").Parse(
		source,
	)

	if err != nil {
		panic(err.Error())
	}

	cleanup := t.Cleanup.String()
	if len(cleanup) > 0 {
		cleanup = fmt.Sprintf(`	s.T().Cleanup(func() {
		%v
	})`, cleanup)
	}

	var result = new(strings.Builder)

	_ = tmpl.Execute(result, struct {
		Dir     string
		Name    string
		Cleanup string
		Run     string
	}{
		Name:    t.Name,
		Dir:     t.Dir,
		Cleanup: cleanup,
		Run:     t.Run.String(),
	})

	return result.String()
}
