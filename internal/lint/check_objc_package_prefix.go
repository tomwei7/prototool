package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/emicklei/proto"
	"github.com/uber/prototool/internal/text"
)

var defaultObjcPrefix = "BAPI"

var fileOptionsObjcPackagePrefixLinter = NewLinter(
	"FILE_OPTIONS_OBJC_PACKAGE_PREFIX",
	fmt.Sprintf(`Verifies that the file option "objc_class_prefix" has prefix "%s", prefix value can set by environment PROTO_OBJC_PACKAGE_PREFIX`, defaultObjcPrefix),
	checkFileOptionsObjcPackagePrefix,
)

func checkFileOptionsObjcPackagePrefix(add func(*text.Failure), dirPath string, descriptors []*FileDescriptor) error {
	return runVisitor(&fileOptionsObjcPackagePrefixVisitor{baseAddVisitor: newBaseAddVisitor(add)}, descriptors)
}

type fileOptionsObjcPackagePrefixVisitor struct {
	baseAddVisitor

	option   *proto.Option
	fileName string
}

func (v *fileOptionsObjcPackagePrefixVisitor) OnStart(descriptor *FileDescriptor) error {
	v.fileName = descriptor.Filename
	v.option = nil
	return nil
}

func (v *fileOptionsObjcPackagePrefixVisitor) VisitOption(element *proto.Option) {
	if element.Name == "objc_class_prefix" {
		v.option = element
	}
}

func (v *fileOptionsObjcPackagePrefixVisitor) Finally() error {
	if v.option == nil {
		return nil
	}
	value := v.option.Constant.Source
	prefix := defaultObjcPrefix
	if v := os.Getenv("PROTO_OBJC_PACKAGE_PREFIX"); v != "" {
		prefix = v
	}
	words := strings.Split(filepath.Dir(v.fileName), "/")
	for _, v := range words {
		prefix += objcUCFirst(v)
	}
	if value != prefix {
		v.AddFailuref(v.option.Position, `Expect option "objc_package" as: "%s" actual: "%s"`, prefix, value)
	}
	return nil
}

func objcUCFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}
