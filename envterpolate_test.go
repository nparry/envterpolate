package main

import (
	"bytes"
	"github.com/bmizerany/assert"
	"testing"
)

var theWordIsGo = map[string]string{
	"WORD": "go",
}

// Wrap substituteVariableReferences for easy testing
func subst(input string, vars map[string]string) string {
	buf := new(bytes.Buffer)
	substituteVariableReferences(bytes.NewBufferString(input), buf, func(s string) string {
		return vars[s]
	})
	return buf.String()
}

func TestEmptyInput(t *testing.T) {
	result := subst("", theWordIsGo)
	assert.Equal(t, result, "")
}

func TestNoVariables(t *testing.T) {
	result := subst("hello world", theWordIsGo)
	assert.Equal(t, result, "hello world")
}

func TestSimpleVariable(t *testing.T) {
	result := subst("hello $WORD world", theWordIsGo)
	assert.Equal(t, result, "hello go world")
}

func TestSimpleVariableAtStart(t *testing.T) {
	result := subst("$WORD home world", theWordIsGo)
	assert.Equal(t, result, "go home world")
}

func TestSimpleVariableAtEnd(t *testing.T) {
	result := subst("let's $WORD", theWordIsGo)
	assert.Equal(t, result, "let's go")
}

func TestOnlyVariable(t *testing.T) {
	result := subst("$WORD", theWordIsGo)
	assert.Equal(t, result, "go")
}

func TestRunOnVariable(t *testing.T) {
	result := subst("$WORD$WORD$WORD!", theWordIsGo)
	assert.Equal(t, result, "gogogo!")
}

func TestRunOnVariableWithNonVariableTextPrefix(t *testing.T) {
	result := subst("$WORD,no$WORD", theWordIsGo)
	assert.Equal(t, result, "go,nogo")
}

func TestSimpleStandAloneDollar(t *testing.T) {
	result := subst("2 $ for your $WORD thoughts", theWordIsGo)
	assert.Equal(t, result, "2 $ for your go thoughts")
}

func TestSimpleStandAloneDollarAtStart(t *testing.T) {
	result := subst("$ for your $WORD thoughts", theWordIsGo)
	assert.Equal(t, result, "$ for your go thoughts")
}

func TestSimpleStandAloneDollarAtEnd(t *testing.T) {
	result := subst("$WORD, find some $", theWordIsGo)
	assert.Equal(t, result, "go, find some $")
}

func TestOnlyStandAloneDollar(t *testing.T) {
	result := subst("$", theWordIsGo)
	assert.Equal(t, result, "$")
}

func TestStandAloneDollarSuffix(t *testing.T) {
	result := subst("$WORD$", theWordIsGo)
	assert.Equal(t, result, "go$")
}
