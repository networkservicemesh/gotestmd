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

package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/gotestmd/internal/config"
	"github.com/networkservicemesh/gotestmd/internal/generator"
	"github.com/networkservicemesh/gotestmd/internal/linker"
	"github.com/networkservicemesh/gotestmd/internal/parser"
)

func main() {
	c := config.FromArgs()
	_ = os.MkdirAll(c.OutputDir, os.ModePerm)
	var examples []*parser.Example

	var p = parser.New()
	var l = linker.New(c.InputDir)
	var g = generator.New(c)
	dirs := getRecursiveDirectories(c.InputDir)
	for _, dir := range dirs {
		ex, err := p.ParseFile(path.Join(dir, "README.md"))
		if err == nil {
			examples = append(examples, ex)
		}
	}
	linkedExamples, err := l.Link(examples...)
	if err != nil {
		logrus.Fatalf("cannot build examples: %v", err.Error())
	}

	suites := g.Generate(linkedExamples...)
	for _, suite := range suites {
		dir, _ := filepath.Split(suite.Location)
		_ = os.MkdirAll(dir, os.ModePerm)
		err := ioutil.WriteFile(suite.Location, []byte(suite.String()), os.ModePerm)
		if err != nil {
			logrus.Fatalf("cannot save suite %v, : %v", suite.Name(), err.Error())
		}
	}
}

func getFilter(root string) func(string) bool {
	var ignored []string
	ignored = append(ignored, filepath.Join(root, ".git"))

	return func(s string) bool {
		for _, line := range ignored {
			if strings.HasPrefix(s, line) {
				return true
			}
		}
		return false
	}
}

func getRecursiveDirectories(root string) []string {
	var result []string
	var isIgnored = getFilter(root)
	_ = filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && !isIgnored(path) {
				result = append(result, path)
			}
			return nil
		})

	return result
}
