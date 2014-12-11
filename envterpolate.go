/*
  Copyright 2014 Nathan Parry

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

// envterpolate - an envsubst semi-clone in go
package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"unicode"
)

type runeReader interface {
	ReadRune() (rune, int, error)
	UnreadRune() error
}

type runeWriter interface {
	WriteRune(rune) (int, error)
}

const (
	initial = iota
	readingVarName
)

func isVarNameCharacter(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_'
}

func outputVarValue(buffer *bytes.Buffer, target runeWriter, resolver func(string) string) error {
	varName := buffer.String()
	varValue := resolver(varName)
	for _, varValueChar := range varValue {
		if _, err := target.WriteRune(varValueChar); err != nil {
			return err
		}
	}

	return nil
}

func substituteVariableReferences(source runeReader, target runeWriter, resolver func(string) string) error {
	buffer := new(bytes.Buffer)
	state := initial
	var err error

	for char, size, _ := source.ReadRune(); size != 0; char, size, _ = source.ReadRune() {
		switch state {
		case initial:
			switch {
			case char == '$':
				state = readingVarName
			default:
				_, err = target.WriteRune(char)
			}
		case readingVarName:
			switch {
			case isVarNameCharacter(char):
				buffer.WriteRune(char)
			default:
				source.UnreadRune()
				state = initial
				err = outputVarValue(buffer, target, resolver)
				buffer.Reset()
			}
		}

		if err != nil {
			return err
		}
	}

	if state == readingVarName {
		err = outputVarValue(buffer, target, resolver)
	}

	return err
}

func main() {
	out := bufio.NewWriter(os.Stdout)
	if err := substituteVariableReferences(bufio.NewReader(os.Stdin), out, os.Getenv); err != nil {
		log.Fatal(err)
	}
	out.Flush()
}
