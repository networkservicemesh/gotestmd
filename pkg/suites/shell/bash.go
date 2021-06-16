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

package shell

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const (
	bufferSize           = 1 << 16
	finishMessage        = "gotestmd/pkg/suites/shell/Bash.const.finish"
	cmdPrintStatusCode   = `echo -e \\n$?`
	cmdPrintStdoutFinish = `echo ` + finishMessage
	cmdPrintStderrFinish = cmdPrintStdoutFinish + ` >&2`
)

// Bash is api for bash process
type Bash struct {
	Dir       string
	Env       []string
	once      sync.Once
	resources []io.Closer
	ctx       context.Context
	cancel    context.CancelFunc

	cmd *exec.Cmd

	stdin    io.Writer
	stdoutCh chan string
	stderrCh chan string
}

// Close closes current bash process and all used resources
func (b *Bash) Close() {
	b.once.Do(b.init)
	b.cancel()
	_, _ = b.stdin.Write([]byte("exit 0\n"))
	_ = b.cmd.Wait()
	for _, r := range b.resources {
		_ = r.Close()
	}
}

func (b *Bash) init() {
	b.ctx, b.cancel = context.WithCancel(context.Background())
	b.stdoutCh = make(chan string)
	b.stderrCh = make(chan string)
	p, err := exec.LookPath("bash")
	if err != nil {
		panic(err.Error())
	}
	if len(b.Env) == 0 {
		b.Env = os.Environ()
	}
	b.cmd = &exec.Cmd{
		Dir:  b.Dir,
		Env:  b.Env,
		Path: p,
	}

	stderr, err := b.cmd.StderrPipe()
	if err != nil {
		panic(err.Error())
	}
	b.resources = append(b.resources, stderr)

	stdin, err := b.cmd.StdinPipe()
	if err != nil {
		panic(err.Error())
	}
	b.resources = append(b.resources, stdin)
	b.stdin = stdin

	stdout, err := b.cmd.StdoutPipe()
	if err != nil {
		panic(err.Error())
	}
	b.resources = append(b.resources, stdout)

	err = b.cmd.Start()
	if err != nil {
		panic(err.Error())
	}

	go b.extractMessagesFromPipe(stdout, b.stdoutCh)
	go b.extractMessagesFromPipe(stderr, b.stderrCh)
}

func (b *Bash) extractMessagesFromPipe(pipe io.Reader, ch chan string) {
	var output string
	var buffer []byte = make([]byte, bufferSize)
	cur := 0
	for b.ctx.Err() == nil {
		n, err := pipe.Read(buffer[cur:])
		if err != nil {
			return
		}
		r := strings.TrimSpace(string(buffer[:cur+n]))
		if strings.HasSuffix(r, finishMessage) {
			if len(r) >= len("\n"+finishMessage) {
				output = r[:len(r)-len("\n"+finishMessage)]
			}
			ch <- output
			output = ""
			cur = 0
			continue
		}
		cur += n
		if cur == bufferSize {
			cur = 0
		}
	}
}

// Run runs the cmd. Returns stdout and stderr as a result.
func (b *Bash) Run(s string) (stdout, stderr string, exitCode int, err error) {
	b.once.Do(b.init)

	if b.ctx.Err() != nil {
		err = b.ctx.Err()
		return
	}

	_, err = b.stdin.Write([]byte(s + "\n"))
	if err != nil {
		return
	}

	_, err = b.stdin.Write([]byte(cmdPrintStatusCode + "\n"))
	if err != nil {
		return
	}

	_, err = b.stdin.Write([]byte(cmdPrintStdoutFinish + "\n"))
	if err != nil {
		return
	}

	_, err = b.stdin.Write([]byte(cmdPrintStderrFinish + "\n"))
	if err != nil {
		return
	}

	select {
	case stdout = <-b.stdoutCh:
	case <-b.ctx.Done():
		err = b.ctx.Err()
		return
	}

	select {
	case stderr = <-b.stderrCh:
	case <-b.ctx.Done():
		err = b.ctx.Err()
		return
	}

	lastLineBreak := strings.LastIndex(stdout, "\n")
	var exitCodeString string
	if lastLineBreak == -1 {
		exitCodeString = stdout
		stdout = ""
	} else {
		exitCodeString = stdout[(lastLineBreak + 1):]
		stdout = strings.TrimSpace(stdout[:lastLineBreak])
	}
	var exitCode64 int64
	exitCode64, err = strconv.ParseInt(exitCodeString, 0, 9)
	if err != nil {
		return
	}
	exitCode = int(exitCode64)

	return
}
