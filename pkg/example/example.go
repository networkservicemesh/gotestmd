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

// Package example provides models and functions for parsing and building examples
package example

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var linkRegex = regexp.MustCompile(`\[.*\]\(.*\)`)

// Example represents a markdown example. Contains all needed for generating suites content.
type Example struct {
	Includes           []string
	Requires           []string
	Run                []string
	Cleanup            []string
	Dir                string
	Name               string
	Childs             []*Example
	Parents            []*Example
	parentDependencies map[string]struct{}
}

func (e *Example) getParentDependencies() map[string]struct{} {
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
	}
	e.parentDependencies = result
	return result
}

// IsLeaf returns true if the example have not children and is not using as a dependency
func (e *Example) IsLeaf() bool {
	return len(e.Childs) == 0 && len(e.Requires) == 0 && len(e.Parents) > 0
}

// Dependencies returns unique dependencies for this example
func (e *Example) Dependencies() []string {
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

func (e *Example) normalize(root string) {
	if len(root) >= len(e.Dir) {
		e.Name = ""
	} else {
		e.Name = filepath.Clean(e.Dir[len(root)+len(string(filepath.Separator)):])
	}

	for i := 0; i < len(e.Includes); i++ {
		e.Includes[i] = filepath.Join(e.Name, e.Includes[i])
	}
	for i := 0; i < len(e.Requires); i++ {
		e.Requires[i] = filepath.Join(e.Name, e.Requires[i])
	}
}

// Build establishes links between examples
func Build(root string, examples []*Example) error {
	index := map[string]*Example{}
	for _, example := range examples {
		example.normalize(root)
		index[example.Name] = example
	}
	for _, example := range examples {
		for _, include := range example.Includes {
			child := index[include]
			if child == nil {
				return errors.Errorf("unknown include %v for example %v", include, example.Name)
			}
			child.Parents = append(child.Parents, example)
			example.Childs = append(example.Childs, child)
		}
	}
	return nil
}

// IsEmpty returns true if exampel have not section Run
func (e *Example) IsEmpty() bool {
	return len(e.Run) == 0
}

func parseSection(section, s string) string {
	start := strings.Index(s, section)
	if start == -1 {
		return ""
	}
	const end = "#"
	s = s[start+len(section):]
	e := strings.Index(s, end)
	if e == -1 {
		return s
	}
	return s[:e]
}

// Parse parses example from a passed dir
func Parse(dir string) (*Example, error) {
	const readmemd = "README.md"
	p := filepath.Join(dir, readmemd)
	f, err := os.Open(filepath.Clean(p))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	v, err := ParseReader(f)
	if err != nil {
		return nil, err
	}
	v.Dir = dir
	return v, nil
}

// ParseReader parses examples from the io.Reader
func ParseReader(r io.Reader) (*Example, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	source := string(bytes)

	parseScript := func(s string) []string {
		var r []string
		const begin = "```bash"
		for start := strings.Index(s, begin); start > 0; start = strings.Index(s, begin) {
			end := strings.Index(s[start+len(begin):], "```") + start + len(begin)
			if end < 0 {
				break
			}
			r = append(r, strings.TrimSpace(s[start+len(begin):end]))
			s = s[end+len("```"):]
		}
		return r
	}

	return &Example{
		Cleanup:  parseScript(parseSection("# Cleanup", source)),
		Run:      parseScript(parseSection("# Run", source)),
		Includes: parseLinks(parseSection("# Includes", source)),
		Requires: parseLinks(parseSection("# Requires", source)),
	}, nil
}

func parseLinks(s string) []string {
	var result []string
	links := linkRegex.FindAllString(s, -1)
	for _, link := range links {
		start := strings.IndexRune(link, '(') + len(string("("))
		end := strings.IndexRune(link, ')')
		result = append(result, link[start:end])
	}
	return result
}
