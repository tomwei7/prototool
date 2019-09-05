package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/uber/prototool/internal/text"
)

var defaultJavaPackagePrefix = "com.bilibili.bapis"

var fileOptionsJavaPackagePrefixLinter = NewLinter(
	"FILE_OPTIONS_JAVA_PACKAGE_PREFIX",
	fmt.Sprintf(`Verifies that the file option "go_package" has prefix "%s", prefix value can set by environment PROTO_JAVA_PACKAGE_PREFIX`, defaultJavaPackagePrefix),
	checkFileOptionsJavePackagePrefix,
)

func checkFileOptionsJavePackagePrefix(add func(*text.Failure), dirPath string, descriptors []*FileDescriptor) error {
	return runVisitor(&fileOptionsJavaPackagePrefixVisitor{baseAddVisitor: newBaseAddVisitor(add)}, descriptors)
}

type fileOptionsJavaPackagePrefixVisitor struct {
	baseAddVisitor

	option   *proto.Option
	fileName string
}

func (v *fileOptionsJavaPackagePrefixVisitor) OnStart(descriptor *FileDescriptor) error {
	v.fileName = descriptor.Filename
	v.option = nil
	return nil
}

func (v *fileOptionsJavaPackagePrefixVisitor) VisitOption(element *proto.Option) {
	if element.Name == "java_package" {
		v.option = element
	}
}

func (v *fileOptionsJavaPackagePrefixVisitor) Finally() error {
	if v.option == nil {
		return nil
	}
	value := v.option.Constant.Source
	prefix := defaultJavaPackagePrefix
	if v := os.Getenv("PROTO_JAVA_PACKAGE_PREFIX"); v != "" {
		prefix = v
	}
	ignoredDirs := os.Getenv("PROTO_JAVA_PACKAGE_PREFIX_IGNORED")
	if ignoredDirs == "" {
		ignoredDirs = "bilibili,extension,third_party"
	}
	for _, ignored := range strings.Split(ignoredDirs, ",") {
		if ignored == strings.Split(v.fileName, "/")[0] {
			return nil
		}
	}
	expect_package := prefix + "." + strings.Replace(filepath.Dir(v.fileName), "/", ".", -1)
	expect_package = strings.Replace(expect_package, "interface", "interfaces", -1)
	expect_package = strings.Replace(expect_package, "-", "_", -1)
	if expect_package != value {
		v.AddFailuref(v.option.Position, `Expect option "java_package" as: "%s" actual: "%s"`, expect_package, value)
	}
	return nil
}
