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
	"io"
	"log"
	"os"
)

type runeWriter interface {
	WriteRune(rune) (int, error)
}

func substituteVariableReferences(source io.RuneReader, target runeWriter, resolver func(string) string) error {
	for r, size, _ := source.ReadRune(); size != 0; r, size, _ = source.ReadRune() {
		if _, err := target.WriteRune(r); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	out := bufio.NewWriter(os.Stdout)
	if err := substituteVariableReferences(bufio.NewReader(os.Stdin), out, os.Getenv); err != nil {
		log.Fatal(err)
	}
	out.Flush()
}
