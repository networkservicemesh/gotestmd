// Copyright (c) 2021 Doc.ai and/or its affiliates.
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

package bash

// Option is an option for the Runner
type Option func(bash *Bash)

// WithDir sets the directory where the bash runner will be located
func WithDir(dir string) Option {
	return func(bash *Bash) {
		bash.dir = dir
	}
}

// WithEnv sets env variables for the bash runner
func WithEnv(env []string) Option {
	return func(bash *Bash) {
		bash.env = env
	}
}
