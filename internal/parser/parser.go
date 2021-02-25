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

// Package parser provides a markdown file reader and model
package parser

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Parser is markdown file reader
type Parser struct {
	linkRegex *regexp.Regexp
}

// New creates new Parser instance
func New() *Parser {
	return &Parser{
		linkRegex: regexp.MustCompile(`\[.*\]\(.*\)`),
	}
}

// ParseFile reads file
func (p *Parser) ParseFile(filePath string) (*Example, error) {
	f, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	v, err := p.Parse(f)
	if err != nil {
		return nil, err
	}
	v.Dir = filepath.Dir(filePath)
	return v, nil
}

// Parse reads io.Reader
func (p *Parser) Parse(r io.Reader) (*Example, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	source := string(bytes)

	parseScript := func(s string) []string {
		const (
			scriptBegin = "```bash"
			scriptEnd   = "```"
		)

		var r []string
		for start := strings.Index(s, scriptBegin); start >= 0; start = strings.Index(s, scriptBegin) {
			start += len(scriptBegin)

			end := strings.Index(s[start:], scriptEnd)
			if end < 0 {
				break
			}
			end += start

			r = append(r, strings.TrimSpace(s[start:end]))
			s = s[end+len(scriptEnd):]
		}
		return r
	}

	return &Example{
		Cleanup:  parseScript(parseSection("# Cleanup", source)),
		Run:      parseScript(parseSection("# Run", source)),
		Includes: p.parseLinks(parseSection("# Includes", source)),
		Requires: p.parseLinks(parseSection("# Requires", source)),
	}, nil
}

func (p *Parser) parseLinks(s string) []string {
	var result []string
	links := p.linkRegex.FindAllString(s, -1)
	for _, link := range links {
		start := strings.IndexRune(link, '(') + len("(")
		end := strings.IndexRune(link, ')')
		result = append(result, link[start:end])
	}
	return result
}

func parseSection(section, s string) string {
	const (
		sectionEnd = "#"
		blockDelim = "```"
	)

	start := strings.Index(s, section)
	if start == -1 {
		return ""
	}

	s = s[start+len(section):]

	end := -1
	for blockEnd := 0; blockEnd > end; {
		if end = strings.Index(s[blockEnd:], sectionEnd); end == -1 {
			return s
		}
		end += blockEnd

		for {
			blockStart := strings.Index(s[blockEnd:], blockDelim)
			if blockStart < 0 {
				break
			}
			blockStart += blockEnd + len(blockDelim)

			if blockStart > end {
				break
			}

			if blockEnd = strings.Index(s[blockStart:], blockDelim); blockEnd < 0 {
				return s
			}
			blockEnd += blockStart + len(blockDelim)
		}
	}
	return s[:end]
}
