/*
Copyright 2019 Cornelius Weig

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"bytes"
	"io"
	"testing"

	"github.com/corneliusweig/rakkess/pkg/rakkess/client"
	"github.com/stretchr/testify/assert"
)

type accessResult map[string]int

func buildAccess() accessResult {
	return make(map[string]int)
}
func (a accessResult) withResult(result int, verbs ...string) accessResult {
	for _, v := range verbs {
		a[v] = result
	}
	return a
}
func (a accessResult) allowed(verbs ...string) accessResult {
	return a.withResult(client.AccessAllowed, verbs...)
}
func (a accessResult) denied(verbs ...string) accessResult {
	return a.withResult(client.AccessDenied, verbs...)
}
func (a accessResult) get() map[string]int {
	return a
}

const HEADER = "NAME       GET  LIST\n"

func TestPrintResults(t *testing.T) {
	tests := []struct {
		name          string
		verbs         []string
		given         []client.Result
		expected      string
		expectedColor string
	}{
		{
			"single result, all allowed",
			[]string{"get", "list"},
			[]client.Result{
				{Name: "resource1", Access: buildAccess().allowed("get", "list").get()},
			},
			HEADER + "resource1  ✔    ✔\n",
			HEADER + "resource1  \033[32m✔\033[0m    \033[32m✔\033[0m\n",
		},
		{
			"single result, all forbidden",
			[]string{"get", "list"},
			[]client.Result{
				{Name: "resource1", Access: buildAccess().denied("get", "list").get()},
			},
			HEADER + "resource1  ✖    ✖\n",
			HEADER + "resource1  \033[31m✖\033[0m    \033[31m✖\033[0m\n",
		},
		{
			"single result, all not applicable",
			[]string{"get", "list"},
			[]client.Result{
				{Name: "resource1", Access: buildAccess().withResult(client.AccessNotApplicable, "get", "list").get()},
			},
			HEADER + "resource1       \n",
			HEADER + "resource1  \033[0m\033[0m     \033[0m\033[0m\n",
		},
		{
			"single result, all ERR",
			[]string{"get", "list"},
			[]client.Result{
				{Name: "resource1", Access: buildAccess().withResult(client.AccessRequestErr, "get", "list").get()},
			},
			HEADER + "resource1  ERR  ERR\n",
			HEADER + "resource1  \033[35mERR\033[0m  \033[35mERR\033[0m\n",
		},
		{
			"single result, mixed",
			[]string{"get", "list"},
			[]client.Result{
				{Name: "resource1", Access: buildAccess().allowed("list").denied("get").get()},
			},
			HEADER + "resource1  ✖    ✔\n",
			"",
		},
		{
			"many results",
			[]string{"get"},
			[]client.Result{
				{Name: "resource1", Access: buildAccess().denied("get").get()},
				{Name: "resource2", Access: buildAccess().allowed("get").get()},
				{Name: "resource3", Access: buildAccess().denied("get").get()},
			},
			"NAME       GET\nresource1  ✖\nresource2  ✔\nresource3  ✖\n",
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			PrintResults(buf, test.verbs, test.given)

			assert.Equal(t, test.expected, buf.String())
		})
	}

	for _, test := range tests[0:4] {
		isTerminal := IsTerminal
		IsTerminal = func(w io.Writer) bool {
			return true
		}
		defer func() {
			IsTerminal = isTerminal
		}()

		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			PrintResults(buf, test.verbs, test.given)

			assert.Equal(t, test.expectedColor, buf.String())
		})
	}
}
