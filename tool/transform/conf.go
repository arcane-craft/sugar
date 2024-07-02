package main

import (
	"github.com/arcane-craft/sugar/tool/transform/exception"
	"github.com/arcane-craft/sugar/tool/transform/lib"
	"github.com/arcane-craft/sugar/tool/transform/question"
	"github.com/arcane-craft/sugar/tool/transform/tryfunc"
)

const (
	SyntaxQuestionMark = "question_mark"
	SyntaxException    = "exception"
	SyntaxTryFunc      = "try_func"
)

var programsByName = map[string]lib.Program{
	SyntaxQuestionMark: new(question.Translator),
	SyntaxException:    new(exception.Translator),
	SyntaxTryFunc:      new(tryfunc.Translator),
}

func Programs(names []string) []lib.Program {
	var progs []lib.Program
	for _, n := range names {
		p, ok := programsByName[n]
		if ok {
			progs = append(progs, p)
		}
	}
	return progs
}

func DefaultSyntax() []string {
	return []string{
		SyntaxQuestionMark,
		SyntaxException,
		SyntaxTryFunc,
	}
}
