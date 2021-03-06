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
	"path"
	"path/filepath"
	"strings"

	"github.com/networkservicemesh/gotestmd/internal/config"
	"github.com/networkservicemesh/gotestmd/internal/linker"
)

// Generator can generate suites from the slice of linker.LinedExample
type Generator struct {
	conf config.Config
}

// New creates new Generator instance
func New(conf config.Config) *Generator {
	return &Generator{
		conf: conf,
	}
}

// Generate generates suites based on passed examples
func (g *Generator) Generate(examples ...*linker.LinkedExample) []*Suite {
	var result []*Suite
	var tests = map[string][]*Test{}
	var index = map[string]*Suite{}
	var children = map[string][]*Suite{}
	moduleName := moduleName(g.conf.OutputDir)
	for _, e := range examples {
		if e.IsLeaf() {
			_, name := path.Split(e.Name)
			for _, parent := range e.Parents {
				tests[parent.Name] = append(tests[parent.Name], &Test{
					Dir:     e.Dir,
					Name:    strings.Title(nameRegex.ReplaceAllString(name, "_")),
					Cleanup: e.Cleanup,
					Run:     e.Run,
				})
			}
			continue
		}

		// Dependencies to import
		var deps = Dependencies([]Dependency{Dependency(g.conf.BasePkg)})
		deps = append(deps, normalizeDeps(moduleName, e.Dependencies())...)

		// Parent suites to setup first
		var depsToSetup = Dependencies([]Dependency{Dependency(g.conf.BasePkg)})
		depsToSetup = append(depsToSetup, normalizeDeps(moduleName, e.ParentDependencies())...)

		s := &Suite{
			Dir:         e.Dir,
			Location:    filepath.Join(g.conf.OutputDir, strings.ToLower(e.Name), "suite.gen.go"),
			Dependency:  Dependency(path.Join(g.conf.OutputDir, strings.ToLower(e.Name))),
			Cleanup:     e.Cleanup,
			Run:         e.Run,
			Deps:        deps,
			DepsToSetup: depsToSetup,
		}

		// Remember if suite is a subsuite
		for _, parent := range e.Parents {
			children[parent.Name] = append(children[parent.Name], s)
		}
		result = append(result, s)

		index[e.Name] = s
	}

	// Apply tests to the suites
	for k, v := range tests {
		index[k].Tests = append(index[k].Tests, v...)
	}

	// Apply subsuites to the suites
	for k, v := range children {
		index[k].Children = append(index[k].Children, v...)
	}

	return result
}
