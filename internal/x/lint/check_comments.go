package lint

import (
	"strings"
	"text/scanner"

	"github.com/emicklei/proto"
	"github.com/uber/prototool/internal/text"
)

var commentsChecker = NewAddChecker(
	"COMMENTS",
	"varifies comments",
	checkComments,
)

func checkComments(add func(*text.Failure), dirPath string, descriptors []*proto.Proto) error {
	return runVisitor(&commentsVisitor{baseAddVisitor: newBaseAddVisitor(add)}, descriptors)
}

type commentsVisitor struct {
	baseAddVisitor
}

func (v commentsVisitor) VisitMessage(element *proto.Message) {
	if !strings.HasSuffix(element.Name, "Req") && !strings.HasSuffix(element.Name, "Reply") {
		v.checkComments(element.Position, element.Comment)
	}
	for _, child := range element.Elements {
		child.Accept(v)
	}
}

func (v commentsVisitor) VisitService(element *proto.Service) {
	v.checkComments(element.Position, element.Comment)
	for _, child := range element.Elements {
		child.Accept(v)
	}
}

func (v commentsVisitor) VisitNormalField(element *proto.NormalField) {
	v.checkComments(element.Position, element.Comment, element.InlineComment)
	for _, child := range element.Options {
		child.Accept(v)
	}
}

func (v commentsVisitor) VisitEnumField(element *proto.EnumField) {
	v.checkComments(element.Position, element.Comment, element.InlineComment)
	if element.ValueOption != nil {
		element.ValueOption.Accept(v)
	}
}

func (v commentsVisitor) VisitEnum(element *proto.Enum) {
	v.checkComments(element.Position, element.Comment)
	for _, child := range element.Elements {
		child.Accept(v)
	}
}

func (v commentsVisitor) VisitOneof(element *proto.Oneof) {
	v.checkComments(element.Position, element.Comment)
	for _, child := range element.Elements {
		child.Accept(v)
	}
}

func (v commentsVisitor) VisitOneofField(element *proto.OneOfField) {
	v.checkComments(element.Position, element.Comment, element.InlineComment)
	for _, child := range element.Options {
		child.Accept(v)
	}
}

func (v commentsVisitor) VisitRPC(element *proto.RPC) {
	v.checkComments(element.Position, element.Comment, element.InlineComment)
}

func (v commentsVisitor) checkComments(position scanner.Position, comments ...*proto.Comment) {
	for _, comment := range comments {
		if comment != nil {
			return
		}
	}
	v.AddFailuref(position, "miss comment")
}
