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
	return substitute(input, remove, vars)
}

func substitute(input string, undefinedBehavior undefinedVariableBehavior, vars map[string]string) string {
	buf := new(bytes.Buffer)
	substituteVariableReferences(bytes.NewBufferString(input), buf, undefinedBehavior, func(s string) string {
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

	result = subst("hello ${WORD} world", theWordIsGo)
	assert.Equal(t, result, "hello go world")
}

func TestSimpleVariableAtStart(t *testing.T) {
	result := subst("$WORD home world", theWordIsGo)
	assert.Equal(t, result, "go home world")

	result = subst("${WORD} home world", theWordIsGo)
	assert.Equal(t, result, "go home world")
}

func TestSimpleVariableAtEnd(t *testing.T) {
	result := subst("let's $WORD", theWordIsGo)
	assert.Equal(t, result, "let's go")

	result = subst("let's ${WORD}", theWordIsGo)
	assert.Equal(t, result, "let's go")
}

func TestOnlyVariable(t *testing.T) {
	result := subst("$WORD", theWordIsGo)
	assert.Equal(t, result, "go")

	result = subst("${WORD}", theWordIsGo)
	assert.Equal(t, result, "go")
}

func TestRunOnVariable(t *testing.T) {
	result := subst("$WORD$WORD$WORD!", theWordIsGo)
	assert.Equal(t, result, "gogogo!")

	result = subst("${WORD}${WORD}${WORD}!", theWordIsGo)
	assert.Equal(t, result, "gogogo!")
}

func TestRunOnVariableWithNonVariableTextPrefix(t *testing.T) {
	result := subst("$WORD,no$WORD", theWordIsGo)
	assert.Equal(t, result, "go,nogo")

	result = subst("${WORD},no${WORD}", theWordIsGo)
	assert.Equal(t, result, "go,nogo")
}

func TestSimpleStandAloneDollar(t *testing.T) {
	result := subst("2 $ for your $WORD thoughts", theWordIsGo)
	assert.Equal(t, result, "2 $ for your go thoughts")

	result = subst("2 ${} for your $WORD thoughts", theWordIsGo)
	assert.Equal(t, result, "2 ${} for your go thoughts")
}

func TestSimpleStandAloneDollarAtStart(t *testing.T) {
	result := subst("$ for your $WORD thoughts", theWordIsGo)
	assert.Equal(t, result, "$ for your go thoughts")

	result = subst("${} for your $WORD thoughts", theWordIsGo)
	assert.Equal(t, result, "${} for your go thoughts")
}

func TestSimpleStandAloneDollarAtEnd(t *testing.T) {
	result := subst("$WORD, find some $", theWordIsGo)
	assert.Equal(t, result, "go, find some $")

	result = subst("${WORD}, find some ${}", theWordIsGo)
	assert.Equal(t, result, "go, find some ${}")
}

func TestOnlyStandAloneDollar(t *testing.T) {
	result := subst("$", theWordIsGo)
	assert.Equal(t, result, "$")

	result = subst("${}", theWordIsGo)
	assert.Equal(t, result, "${}")
}

func TestStandAloneDollarSuffix(t *testing.T) {
	result := subst("$WORD$", theWordIsGo)
	assert.Equal(t, result, "go$")

	result = subst("${WORD}${}", theWordIsGo)
	assert.Equal(t, result, "go${}")
}

func TestBracingStartsMidtoken(t *testing.T) {
	result := subst("$WORD{FISH}", theWordIsGo)
	assert.Equal(t, result, "go{FISH}")

	result = subst("$WORD{WORD}", theWordIsGo)
	assert.Equal(t, result, "go{WORD}")

	result = subst("$WO{RD}", theWordIsGo)
	assert.Equal(t, result, "{RD}")
}

func TestUnclosedBraceWithPrefix(t *testing.T) {
	result := subst("$WORD{WORD", theWordIsGo)
	assert.Equal(t, result, "go{WORD")
}

func TestUnclosedBracedDollar(t *testing.T) {
	result := subst("${", theWordIsGo)
	assert.Equal(t, result, "${")
}

func TestUnclosedBracedDollarAndSubsequentTokens(t *testing.T) {
	result := subst("${ stuff }", theWordIsGo)
	assert.Equal(t, result, "${ stuff }")
}

func TestUnclosedBracedDollarWithSuffix(t *testing.T) {
	result := subst("${WORD", theWordIsGo)
	assert.Equal(t, result, "${WORD")
}

func TestUnclosedBracedDollarWithSuffixAndSubsequentTokens(t *testing.T) {
	result := subst("no ${WORD yet", theWordIsGo)
	assert.Equal(t, result, "no ${WORD yet")
}

func TestDoubleBracedDollar(t *testing.T) {
	result := subst("this ${{WORD}} should not be touched", theWordIsGo)
	assert.Equal(t, result, "this ${{WORD}} should not be touched")
}

func TestUndefinedVariablesAreRemovedByDefault(t *testing.T) {
	result := subst("nothing $WORD2 to see here", theWordIsGo)
	assert.Equal(t, result, "nothing  to see here")

	result = subst("nothing ${WORD2} to see here", theWordIsGo)
	assert.Equal(t, result, "nothing  to see here")

	result = subst("nothing${WORD2}to see here", theWordIsGo)
	assert.Equal(t, result, "nothingto see here")
}

func TestUndefinedVariablesArePreservedWhenWanted(t *testing.T) {
	result := substitute("nothing $WORD2 to see here", preserve, theWordIsGo)
	assert.Equal(t, result, "nothing $WORD2 to see here")

	result = substitute("nothing ${WORD2} to see here", preserve, theWordIsGo)
	assert.Equal(t, result, "nothing ${WORD2} to see here")

	result = substitute("nothing${WORD2}to see here", preserve, theWordIsGo)
	assert.Equal(t, result, "nothing${WORD2}to see here")
}
