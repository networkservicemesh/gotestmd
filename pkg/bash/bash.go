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

package bash

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	initialBufferSize    = 1 << 10
	finishMessage        = "gotestmd/pkg/suites/shell/Bash.const.finish"
	cmdPrintStatusCode   = `echo -e \\n$?`
	cmdPrintStdoutFinish = `echo ` + finishMessage
	cmdPrintStderrFinish = cmdPrintStdoutFinish + ` >&2`
)

// Bash is api for bash process
type Bash struct {
	env       []string
	resources []io.Closer
	ctx       context.Context
	cancel    context.CancelFunc

	cmd *exec.Cmd

	stdin    io.Writer
	stdoutCh chan string
	stderrCh chan string
}

// New creates a new bash runner and initializes it
func New(options ...Option) (*Bash, error) {
	b := &Bash{}
	for _, o := range options {
		o(b)
	}

	err := b.Init()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Close closes current bash process and all the resources used by it
func (b *Bash) Close() {
	b.cancel()
	_, err := b.stdin.Write([]byte("exit 0\n"))
	if err != nil {
		panic(err)
	}
	_ = b.cmd.Wait()
	for _, r := range b.resources {
		_ = r.Close()
	}
}

// Init initializes all resources for the bash runner.
//
// Must be called before any call to Run or Close
//
// You are advised to use bash.New instead, which calls this function automatically.
func (b *Bash) Init() error {
	b.ctx, b.cancel = context.WithCancel(context.Background())
	b.stdoutCh = make(chan string)
	b.stderrCh = make(chan string)
	p, err := exec.LookPath("bash")
	if err != nil {
		return err
	}
	if len(b.env) == 0 {
		b.env = os.Environ()
	}
	b.cmd = &exec.Cmd{
		Env:  b.env,
		Path: p,
	}

	stderr, err := b.cmd.StderrPipe()
	if err != nil {
		return err
	}
	b.resources = append(b.resources, stderr)

	stdin, err := b.cmd.StdinPipe()
	if err != nil {
		return err
	}
	b.resources = append(b.resources, stdin)
	b.stdin = stdin

	stdout, err := b.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	b.resources = append(b.resources, stdout)

	err = b.cmd.Start()
	if err != nil {
		return err
	}

	go b.extractMessagesFromPipe(stdout, b.stdoutCh)
	go b.extractMessagesFromPipe(stderr, b.stderrCh)

	return nil
}

func (b *Bash) extractMessagesFromPipe(pipe io.Reader, ch chan string) {
	var buffer = make([]byte, initialBufferSize)
	cur := 0
	for b.ctx.Err() == nil {
		n, err := pipe.Read(buffer[cur:])
		if err != nil {
			return
		}
		cur += n
		r := strings.TrimSpace(string(buffer[:cur]))
		if strings.HasSuffix(r, finishMessage) {
			if len(r) >= len(finishMessage) {
				r = strings.TrimSpace(r[:len(r)-len(finishMessage)])
			}
			select {
			case ch <- r:
			case <-b.ctx.Done():
				return
			}
			cur = 0
			continue
		}
		if cur == len(buffer) {
			oldBuffer := buffer
			buffer = make([]byte, len(oldBuffer)*2)
			copy(buffer, oldBuffer)
		}
	}
}

// Run runs the command
func (b *Bash) Run(cmd string) (stdout, stderr string, exitCode int, err error) {
	if b.ctx.Err() != nil {
		err = b.ctx.Err()
		return
	}

	_, err = b.stdin.Write([]byte(cmd + "\n" + cmdPrintStatusCode + "\n" + cmdPrintStdoutFinish + "\n" + cmdPrintStderrFinish + "\n"))
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

	return stdout, stderr, exitCode, nil
}
