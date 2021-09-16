package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/spy16/slurp/core"
	"github.com/stretchr/testify/assert"
)

func readEval(t *testing.T, lisp string) (core.Any, error) {
	sl := NewInterpreter()
	srf := NewScriptReaderFactory()
	br := bufio.NewReader(strings.NewReader(lisp))
	return srf.ReadEval(sl, br)
}

func TestMakePub(t *testing.T) {
	r, err := readEval(t, "(def out (pub \"tcp://127.0.0.1:9991\"))")
	assert.NoError(t, err, "make pub stream")
	assert.NotNil(t, r, "return a value")
}

func TestMakeSub(t *testing.T) {
	r, err := readEval(t, "(def in (sub \"tcp://127.0.0.1:9991\"))")
	assert.NoError(t, err, "make sub stream")
	assert.NotNil(t, r, "return a value")
}

func TestPrint(t *testing.T) {
	r, err := readEval(t, "(print \"test\")")
	assert.NoError(t, err, "print something")
	assert.NotNil(t, r, "return a value")
	assert.Equal(t, r, "nil")
}

func TestSheBang(t *testing.T) {
	r, err := readEval(t, "#!foo\n1")
	assert.NoError(t, err, "nop after she-bang")
	assert.NotNil(t, r, "return a value")
	assert.Equal(t, r, "1")
}

func writeLisp(lisp string) *os.File {
	tmpf, err := ioutil.TempFile("", "tq")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpf.Name()) // clean up

	if _, err := tmpf.Write([]byte(lisp)); err != nil {
		log.Fatal(err)
	}
	if err := tmpf.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpf
}

func WriteEvalPrint(t *testing.T) {
	tmpf := writeLisp(`#!foo\\
(print test)
1
	`)
	sl := NewInterpreter()
	r := readEvalFiles(sl, tmpf.Name())
	assert.NotNil(t, r, "return a value")
	assert.Equal(t, r, "1")
}
