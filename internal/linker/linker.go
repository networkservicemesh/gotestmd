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

// Package linker provides a linker to add links between examples
package linker

import (
	"github.com/pkg/errors"

	"github.com/networkservicemesh/gotestmd/internal/parser"
)

// Linker can add links between examples
type Linker struct {
	root string
}

// New creates new Linker instance
func New(root string) *Linker {
	return &Linker{
		root: root,
	}
}

// Link adds all possible links between examples. Return error if any link is invalid
func (l *Linker) Link(examples ...*parser.Example) ([]*LinkedExample, error) {
	index := map[string]*LinkedExample{}
	var result []*LinkedExample
	for _, example := range examples {
		linkedExample := NewLinkedExample(l.root, example)
		index[linkedExample.Name] = linkedExample
		result = append(result, linkedExample)
	}
	for _, linkedExample := range result {
		for _, include := range linkedExample.Includes {
			child := index[include]
			if child == nil {
				return nil, errors.Errorf("unknown include %v for example %v", include, linkedExample.Name)
			}
			child.Parents = append(child.Parents, linkedExample)
			linkedExample.Childs = append(linkedExample.Childs, child)
		}
	}
	for _, linkedExample := range result {
		var filteredRequires []string
		for _, require := range linkedExample.Requires {
			if _, ok := linkedExample.getParentDependencies()[require]; !ok {
				filteredRequires = append(filteredRequires, require)
			}
		}
		linkedExample.Requires = filteredRequires
	}
	return result, nil
}
