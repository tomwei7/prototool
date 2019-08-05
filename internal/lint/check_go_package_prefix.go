package lint

import (
	"fmt"
	"os"
	"strings"

	"github.com/emicklei/proto"
	"github.com/uber/prototool/internal/text"
)

var defaultPrefix = "git.bilibili.co/bapis/bapis-go"

var fileOptionsGoPackagePrefixLinter = NewLinter(
	"FILE_OPTIONS_GO_PACKAGE_PREFIX",
	fmt.Sprintf(`Verifies that the file option "go_package" has prefix "%s", prefix value can set by environment PROTO_GO_PACKAGE_PREFIX`, defaultPrefix),
	checkFileOptionsGoPackagePrefix,
)

func checkFileOptionsGoPackagePrefix(add func(*text.Failure), dirPath string, descriptors []*FileDescriptor) error {
	return runVisitor(&fileOptionsGoPackagePrefixVisitor{baseAddVisitor: newBaseAddVisitor(add)}, descriptors)
}

type fileOptionsGoPackagePrefixVisitor struct {
	baseAddVisitor

	option *proto.Option
}

func (v *fileOptionsGoPackagePrefixVisitor) OnStart(descriptor *FileDescriptor) error {
	v.option = nil
	return nil
}

func (v *fileOptionsGoPackagePrefixVisitor) VisitOption(element *proto.Option) {
	if element.Name == "go_package" {
		v.option = element
	}
}

func (v *fileOptionsGoPackagePrefixVisitor) Finally() error {
	if v.option == nil {
		return nil
	}
	value := v.option.Constant.Source
	prefix := defaultPrefix
	if v := os.Getenv("PROTO_GO_PACKAGE_PREFIX"); v != "" {
		prefix = v
	}
	if !strings.HasPrefix(value, prefix) {
		v.AddFailuref(v.option.Position, `Option "go_package" must has prefix: %s`, prefix)
	}
	return nil
}
