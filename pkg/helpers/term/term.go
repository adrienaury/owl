/*
Copyright 2016 The Kubernetes Authors.
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

package term

import (
	"io"
	"os"

	"github.com/docker/docker/pkg/term"
)

// IsTerminal returns whether the passed object is a terminal or not
func IsTerminal(i interface{}) bool {
	_, terminal := term.GetFdInfo(i)
	return terminal
}

// IsTerminalReader returns whether the passed io.Reader is a terminal or not
func IsTerminalReader(r io.Reader) bool {
	file, ok := r.(*os.File)
	return ok && term.IsTerminal(file.Fd())
}

// IsTerminalWriter returns whether the passed io.Writer is a terminal or not
func IsTerminalWriter(w io.Writer) bool {
	file, ok := w.(*os.File)
	return ok && term.IsTerminal(file.Fd())
}
