package main

import (
	"bytes"
	"github.com/bmizerany/assert"
	"testing"
)

// Wrap substituteVariableReferences for easy testing
func subst(input string, vars map[string]string) string {
	buf := new(bytes.Buffer)
	substituteVariableReferences(bytes.NewBufferString(input), buf, func(s string) string {
		return vars[s]
	})
	return buf.String()
}

func TestNoVariables(t *testing.T) {
	result := subst("hello world", map[string]string{
		"DUMMY": "dummy",
	})
	assert.Equal(t, result, "hello world")
}

func TestSimpleVariable(t *testing.T) {
	result := subst("hello $WORD world", map[string]string{
		"WORD": "go",
	})
	assert.Equal(t, result, "hello go world")
}
