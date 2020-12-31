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
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/networkservicemesh/gotestmd/pkg/config"
	"github.com/networkservicemesh/gotestmd/pkg/example"
)

var nameRegex = regexp.MustCompile("[^a-zA-Z0-9]+")
var spaceRegex = regexp.MustCompile(`[\t\r\n]+`)

func normalizeName(s string) string {
	return strings.ToLower(nameRegex.ReplaceAllString(s, "_"))
}

func normalizeDeps(rootPackage string, deps []string) Dependencies {
	var d Dependencies
	for _, dep := range deps {
		pieces := strings.Split(filepath.Clean(dep), string(filepath.Separator))
		for i := 0; i < len(pieces); i++ {
			pieces[i] = normalizeName(pieces[i])
		}
		dep = path.Join(pieces...)
		d = append(d, Dependency(path.Join(rootPackage, dep)))
	}
	return d
}

// Generate generates suites based on passed examples
func Generate(c config.Config, examples []*example.Example) []*SuiteTemplate {
	var result []*SuiteTemplate
	var tests = map[string][]*TestTemplate{}
	var index = map[string]*SuiteTemplate{}

	for _, e := range examples {
		dir := filepath.Clean(strings.TrimPrefix(e.Dir, os.Getenv("GOPATH")))
		if e.IsLeaf() {
			_, name := path.Split(e.Name)
			for _, parent := range e.Parents {
				tests[parent.Name] = append(tests[parent.Name], &TestTemplate{
					Dir:     dir,
					Name:    strings.Title(nameRegex.ReplaceAllString(name, "_")),
					Cleanup: e.Cleanup,
					Run:     e.Run,
				})
			}
			continue
		}
		result = append(result, &SuiteTemplate{
			Dir:        dir,
			Location:   filepath.Join(c.OutputDir, strings.ToLower(e.Name), "suite.gen.go"),
			Dependency: Dependency(path.Join(c.RootPackage, strings.ToLower(e.Name))),
			Cleanup:    e.Cleanup,
			Run:        e.Run,
			Deps:       normalizeDeps(c.RootPackage, e.Dependencies()),
		})
		index[e.Name] = result[len(result)-1]
	}

	for k, v := range tests {
		index[k].Tests = append(index[k].Tests, v...)
	}

	return result
}
