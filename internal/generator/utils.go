// Copyright (c) 2020-2021 Doc.ai and/or its affiliates.
//
// Copyright (c) 2023 Cisco and/or its affiliates.
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
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var nameRegex = regexp.MustCompile("[^a-zA-Z0-9]+")
var spaceRegex = regexp.MustCompile(`[\t\r\n]+`)

func normalizeName(s string) string {
	return strings.ToLower(nameRegex.ReplaceAllString(s, "_"))
}

func normalizeDeps(module string, deps []string) Dependencies {
	var d Dependencies
	for _, dep := range deps {
		pieces := strings.Split(filepath.Clean(dep), string(filepath.Separator))
		for i := 0; i < len(pieces); i++ {
			pieces[i] = normalizeName(pieces[i])
		}
		pieces = append([]string{module}, pieces...)
		d = append(d, Dependency(filepath.Join(pieces...)))
	}
	return d
}

func moduleName(start string) string {
	const gomod = "go.mod"
	currDir, err := filepath.Abs(start)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	for len(currDir) > 0 {
		p := filepath.Clean(filepath.Join(currDir, "go.mod"))
		if _, err = os.Open(p); err == nil {
			source, err := os.ReadFile(p)
			if err != nil {
				logrus.Fatal(err.Error())
			}
			moduleName := strings.TrimPrefix(strings.Split(string(source), "\n")[0], "module ")
			return filepath.Clean(filepath.Join(moduleName, start))
		}
		currDir = filepath.Dir(currDir)
	}
	return ""
}
