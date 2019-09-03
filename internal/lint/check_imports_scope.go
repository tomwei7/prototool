// Copyright (c) 2019 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package lint

import (
	"os"
	"strings"

	"github.com/emicklei/proto"
	"github.com/uber/prototool/internal/text"
)

var (
	importsScopeLinter = NewLinter(
		"IMPORTS_SCOPE",
		`Verifies that there are no weak imports.`,
		checkImportsScope,
	)
)

func checkImportsScope(add func(*text.Failure), dirPath string, descriptors []*FileDescriptor) error {
	return runVisitor(&importsScope{baseAddVisitor: newBaseAddVisitor(add)}, descriptors)
}

type importsScope struct {
	baseAddVisitor
	fileName string
}

func (v *importsScope) OnStart(descriptor *FileDescriptor) error {
	v.fileName = descriptor.Filename
	return nil
}

func (v importsScope) VisitImport(element *proto.Import) {
	baseScope := strings.Split(v.fileName, "/")[0]
	importBaseScope := strings.Split(element.Filename, "/")[0]
	for _, ignore := range strings.Split(os.Getenv("IGNORE_DIR")+",google,github.com,extension", ",") {
		if importBaseScope == ignore {
			return
		}
	}
	if baseScope != importBaseScope {
		v.AddFailuref(element.Position, `invalid import can not import %s from %s`, element.Filename, v.fileName)
	}
}
