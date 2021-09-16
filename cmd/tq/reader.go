package main

import (
	"io"

	"github.com/spy16/slurp"
	"github.com/spy16/slurp/core"
	"github.com/spy16/slurp/reader"

)

func skip(r *reader.Reader, init rune) (core.Any, error) {
	for rn, err := r.NextRune(); rn != '\n'; rn, err = r.NextRune() {
		if err != nil {
			return nil, err
		}
	}
	return nil, reader.ErrSkip
}

type ScriptReaderFactory struct{}

func NewScriptReaderFactory() *ScriptReaderFactory {
	srf := ScriptReaderFactory{}
	return &srf
}

func (*ScriptReaderFactory) NewReader(r io.Reader) *reader.Reader {
	rd := reader.New(r)
	rd.SetMacro('!', true, skip)

	return rd
}

func (srf *ScriptReaderFactory) ReadEval(sl *slurp.Interpreter, r io.Reader) (core.Any, error) {
	mod, err := srf.NewReader(r).All()
	if err != nil {
		return nil, err
	}

	return sl.Eval(mod)
}
