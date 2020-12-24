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

// Package config contains gotestmd configuration
package config

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Config contains input dir with .md examples and output dir for generated suites
type Config struct {
	InputDir  string
	OutputDir string
}

// FromArgs returns Config from the os.Args
func FromArgs() Config {
	if len(os.Args) != 3 {
		logrus.Fatal("ARGs have wrong length. Expected: input-dir output-dir")
	}

	inputDir, _ := filepath.Abs(os.Args[1])
	if _, err := os.Open(filepath.Clean(inputDir)); err != nil {
		logrus.Fatalf("An error during checking dir: %v, error: %v", inputDir, err.Error())
	}
	if _, err := os.Open(filepath.Clean(os.Args[2])); err != nil {
		err = os.MkdirAll(filepath.Clean(os.Args[2]), os.ModePerm)
		if err != nil {
			logrus.Fatalf("An error during creating dir: %v, error: %v", os.Args[2], err.Error())
		}
	}
	outputDir, _ := filepath.Abs(os.Args[2])
	return Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
	}
}
