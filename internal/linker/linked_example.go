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

package linker

import (
	"path/filepath"

	"github.com/networkservicemesh/gotestmd/internal/parser"
)

// LinkedExample represents parser.Example with links
type LinkedExample struct {
	*parser.Example
	Name               string
	Childs             []*LinkedExample
	Parents            []*LinkedExample
	parentDependencies map[string]struct{}
}

func (e *LinkedExample) getParentDependencies() map[string]struct{} {
	if e.parentDependencies != nil {
		return e.parentDependencies
	}
	var result = make(map[string]struct{})
	for _, parent := range e.Parents {
		for _, dep := range parent.Dependencies() {
			result[dep] = struct{}{}
		}
		for dep := range parent.getParentDependencies() {
			result[dep] = struct{}{}
		}
		result[parent.Name] = struct{}{}
	}
	e.parentDependencies = result
	return result
}

// IsLeaf returns true if the example have not children and is not using as a dependency
func (e *LinkedExample) IsLeaf() bool {
	return len(e.Childs) == 0 && len(e.Requires) == 0 && len(e.Parents) > 0
}

// Dependencies returns unique dependecies for this example
func (e *LinkedExample) Dependencies() []string {
	var deps []string
	var parentDeps []string

	for _, child := range e.Childs {
		if child.IsLeaf() {
			continue
		}
		deps = append(deps, child.Name)
	}

	for _, dep := range e.Requires {
		if _, ok := e.getParentDependencies()[dep]; !ok {
			parentDeps = append(parentDeps, dep)
		}
	}

	return append(parentDeps, deps...)
}

// NewLinkedExample creates new linked example based on parser.Example
func NewLinkedExample(root string, e *parser.Example) *LinkedExample {
	var result = new(LinkedExample)
	result.Example = e
	if len(root) >= len(e.Dir) {
		result.Name = ""
	} else {
		result.Name = filepath.Clean(e.Dir[len(root):])
	}

	for i := 0; i < len(e.Includes); i++ {
		e.Includes[i] = filepath.Join(result.Name, e.Includes[i])
	}
	for i := 0; i < len(e.Requires); i++ {
		e.Requires[i] = filepath.Join(result.Name, e.Requires[i])
	}

	return result
}
