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
}

type runeWriter interface {
	WriteRune(rune) (int, error)
}

type state int

const (
	initial state = iota
	readingVarName
	readingBracedVarName
)

type varNameTokenStatus int

const (
	complete varNameTokenStatus = iota
	incomplete
)

type envterpolator struct {
	state    state
	buffer   bytes.Buffer
	target   runeWriter
	resolver func(string) string
}

func isVarNameCharacter(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_'
}

func standaloneDollarString(varNameTokenStatus varNameTokenStatus, state state) string {
	switch {
	case state == readingVarName:
		return "$"
	case varNameTokenStatus == incomplete:
		return "${"
	}

	return "${}"
}

func writeString(s string, target runeWriter) error {
	for _, char := range s {
		if err := writeRune(char, target); err != nil {
			return err
		}
	}

	return nil
}

func writeRune(char rune, target runeWriter) error {
	_, err := target.WriteRune(char)
	return err
}

func substituteVariableReferences(source runeReader, target runeWriter, resolver func(string) string) error {
	et := envterpolator{
		target:   target,
		resolver: resolver,
	}

	for char, size, _ := source.ReadRune(); size != 0; char, size, _ = source.ReadRune() {
		if err := et.processRune(char); err != nil {
			return err
		}
	}

	return et.endOfInput()
}

func (et *envterpolator) processRune(char rune) error {
	switch et.state {
	case initial:
		switch {
		case char == '$':
			et.state = readingVarName
		default:
			return writeRune(char, et.target)
		}
	case readingVarName:
		switch {
		case isVarNameCharacter(char):
			return writeRune(char, &et.buffer)
		case char == '{' && et.buffer.Len() == 0:
			et.state = readingBracedVarName
		default:
			return et.flushBufferAndProcessNextRune(complete, char)
		}
	case readingBracedVarName:
		switch {
		case isVarNameCharacter(char):
			return writeRune(char, &et.buffer)
		case char == '}':
			return et.flushBuffer(complete)
		default:
			return et.flushBufferAndProcessNextRune(incomplete, char)
		}
	}

	return nil
}

func (et *envterpolator) endOfInput() error {
	if et.state != initial {
		return et.flushBuffer(incomplete)
	}

	return nil
}

func (et *envterpolator) flushBufferAndProcessNextRune(bufferStatus varNameTokenStatus, nextChar rune) error {
	if err := et.flushBuffer(bufferStatus); err != nil {
		return err
	}

	return et.processRune(nextChar)
}

func (et *envterpolator) flushBuffer(bufferStatus varNameTokenStatus) error {
	var err error

	switch {
	case et.buffer.Len() == 0:
		err = writeString(standaloneDollarString(bufferStatus, et.state), et.target)
	case et.state == readingBracedVarName && bufferStatus == incomplete:
		err = writeString("${"+et.buffer.String(), et.target)
	default:
		err = writeString(et.resolver(et.buffer.String()), et.target)
	}

	et.state = initial
	et.buffer.Reset()

	return err
}

func main() {
	out := bufio.NewWriter(os.Stdout)
	if err := substituteVariableReferences(bufio.NewReader(os.Stdin), out, os.Getenv); err != nil {
		log.Fatal(err)
	}
	out.Flush()
}
