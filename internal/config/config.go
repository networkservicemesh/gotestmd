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

// Package config contains gotestmd configuration
package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Config contains input dir with .md examples and output dir for generated suites
type Config struct {
	InputDir  string
	OutputDir string
	BasePkg   string
}

// FromArgs returns Config from the os.Args
func FromArgs() Config {
	if len(os.Args) < 3 || len(os.Args) > 4 {
		logrus.Fatal("ARGs have wrong length. Expected: (string)input-dir (string)output-dir (string)base-pkg[optional]")
	}
	result := Config{
		InputDir:  os.Args[1],
		OutputDir: os.Args[2],
		BasePkg:   "github.com/networkservicemesh/gotestmd/pkg/suites/shell",
	}

	if len(os.Args) == 4 {
		result.BasePkg = os.Args[3]
	}

	return result
}
