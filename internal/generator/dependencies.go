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
	"path/filepath"
	"strings"
)

// Dependency represents test dependency
type Dependency string

// Pkg returns a string that can be imported
func (d Dependency) Pkg() string {
	return string(d)
}

// Name returns pkg name
func (d Dependency) Name() string {
	_, name := filepath.Split(d.Pkg())
	return normalizeName(name)
}

// Dependencies represent an array of Dependency
type Dependencies []Dependency

// FieldsString returns a string that contains a declaration of suite dependencies as fields
func (d Dependencies) FieldsString() string {
	var result strings.Builder
	for i := 0; i < len(d); i++ {
		if i != 0 {
			_, _ = result.WriteString(d[i].Name())
			_, _ = result.WriteString("Suite ")
		}
		_, _ = result.WriteString(d[i].Name())
		_, _ = result.WriteString(".Suite")
		if i+1 < len(d) {
			_, _ = result.WriteString("\n")
		}
	}

	return result.String()
}

// SetupString returns a string that contains a declaration of suite dependencies as part of setup function
func (d Dependencies) SetupString() string {
	if len(d) == 0 {
		return ""
	}

	var result strings.Builder

	result.WriteString("parents := []interface{}{")
	result.WriteString("&s.Suite")
	for i := 1; i < len(d); i++ {
		if i < len(d) {
			result.WriteString(",")
		}
		result.WriteString("&s.")
		result.WriteString(d[i].Name())
		result.WriteString("Suite")
	}
	result.WriteString("}\n")

	result.WriteString(`for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
`)

	return result.String()
}

// String returns a string that contains a declaration of suite dependencies as part of import
func (d Dependencies) String() string {
	var result strings.Builder

	if len(d) > 0 {
		_, _ = result.WriteString(`"github.com/stretchr/testify/suite"`)
		_, _ = result.WriteString("\n")
	}

	for i := 0; i < len(d); i++ {
		_, _ = result.WriteString("\"")
		_, _ = result.WriteString(d[i].Pkg())
		_, _ = result.WriteString("\"")
		if i+1 < len(d) {
			_, _ = result.WriteString("\n")
		}
	}

	return result.String()
}
