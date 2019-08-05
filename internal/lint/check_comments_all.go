package lint

import (
	"strings"
	"text/scanner"

	"github.com/emicklei/proto"
	"github.com/uber/prototool/internal/text"
)

var commentsAllLinter = NewLinter(
	"COMMENTS_ALL",
	"varifies that every field has comment.",
	checkCommentsAll,
)

func checkCommentsAll(add func(*text.Failure), dirPath string, descriptors []*FileDescriptor) error {
	return runVisitor(&commentsAllVisitor{baseAddVisitor: newBaseAddVisitor(add)}, descriptors)
}

type commentsAllVisitor struct {
	baseAddVisitor
}

func (v commentsAllVisitor) VisitMessage(element *proto.Message) {
	if !strings.HasSuffix(element.Name, "Req") && !strings.HasSuffix(element.Name, "Reply") {
		v.checkComments(element.Position, element.Comment)
	}
	for _, child := range element.Elements {
		child.Accept(v)
	}
}

func (v commentsAllVisitor) VisitService(element *proto.Service) {
	v.checkComments(element.Position, element.Comment)
	for _, child := range element.Elements {
		child.Accept(v)
	}
}

func (v commentsAllVisitor) VisitNormalField(element *proto.NormalField) {
	v.checkComments(element.Position, element.Comment, element.InlineComment)
	for _, child := range element.Options {
		child.Accept(v)
	}
}

func (v commentsAllVisitor) VisitEnumField(element *proto.EnumField) {
	v.checkComments(element.Position, element.Comment, element.InlineComment)
	if element.ValueOption != nil {
		element.ValueOption.Accept(v)
	}
}

func (v commentsAllVisitor) VisitEnum(element *proto.Enum) {
	v.checkComments(element.Position, element.Comment)
	for _, child := range element.Elements {
		child.Accept(v)
	}
}

func (v commentsAllVisitor) VisitOneof(element *proto.Oneof) {
	v.checkComments(element.Position, element.Comment)
	for _, child := range element.Elements {
		child.Accept(v)
	}
}

func (v commentsAllVisitor) VisitOneofField(element *proto.OneOfField) {
	v.checkComments(element.Position, element.Comment, element.InlineComment)
	for _, child := range element.Options {
		child.Accept(v)
	}
}

func (v commentsAllVisitor) VisitRPC(element *proto.RPC) {
	v.checkComments(element.Position, element.Comment, element.InlineComment)
}

func (v commentsAllVisitor) checkComments(position scanner.Position, comments ...*proto.Comment) {
	for _, comment := range comments {
		if comment != nil {
			return
		}
	}
	v.AddFailuref(position, "miss comment")
}
