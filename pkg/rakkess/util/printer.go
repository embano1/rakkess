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
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/corneliusweig/rakkess/pkg/rakkess/client"
	"golang.org/x/crypto/ssh/terminal"
)

type color int

const (
	red    = color(31)
	green  = color(32)
	purple = color(35)
	none   = color(0)
)

var IsTerminal = isTerminal

func PrintResults(out io.Writer, requestedVerbs []string, results []client.Result) {
	w := NewWriter(out, 4, 8, 2, ' ', CollapseEscape^StripEscape)
	defer w.Flush()

	fmt.Fprint(w, "NAME")
	for _, v := range requestedVerbs {
		fmt.Fprintf(w, "\t%s", strings.ToUpper(v))
	}
	fmt.Fprint(w, "\n")

	codeConverter := humanreadableAccessCode
	if IsTerminal(out) {
		codeConverter = colorHumanreadableAccessCode
	}

	for _, r := range results {
		fmt.Fprintf(w, "%s", r.Name)
		for _, v := range requestedVerbs {
			fmt.Fprintf(w, "\t%s", codeConverter(r.Access[v]))
		}
		fmt.Fprint(w, "\n")
	}
}

func humanreadableAccessCode(code int) string {
	switch code {
	case client.AccessAllowed:
		return "✔" // ✓
	case client.AccessDenied:
		return "✖" // ✕
	case client.AccessNotApplicable:
		return ""
	case client.AccessRequestErr:
		return "ERR"
	default:
		panic("unknown access code")
	}
}

func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return terminal.IsTerminal(int(f.Fd()))
	}
	return false
}

func colorHumanreadableAccessCode(code int) string {
	return fmt.Sprintf("\xff\033[%dm\xff%s\xff\033[0m\xff", codeToColor(code), humanreadableAccessCode(code))
}

func codeToColor(code int) color {
	switch code {
	case client.AccessAllowed:
		return green
	case client.AccessDenied:
		return red
	case client.AccessNotApplicable:
		return none
	case client.AccessRequestErr:
		return purple
	}
	return none
}
