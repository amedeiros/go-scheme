package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/amedeiros/go-scheme/object"
)

var TRUE = &object.Boolean{Value: true}
var FALSE = &object.Boolean{Value: false}

type Reader struct {
	reader *bufio.Reader
}

func NewReader(input string) *Reader {
	return &Reader{bufio.NewReader(strings.NewReader(input))}
}

func (r *Reader) Read() object.Object {
	char, err := r.currentByte()
	if err != nil {
		return err
	}

	switch char {
	case '#':
		val, err := r.peek()
		if err != nil {
			return err
		}

		peekChar := strings.ToUpper(string(val))

		if peekChar == "T" {
			r.skip()
			return TRUE
		} else if peekChar == "F" {
			r.skip()
			return FALSE
		} else if peekChar == "\\" {
			r.skip()
			cur, err := r.currentByte()
			if err != nil {
				return err
			}
			return &object.Char{Value: string(cur)}
		}

		panic(fmt.Sprintf("Expecting one of F or T or \\ found %s instead.", peekChar))
	case '"':
		str := bytes.Buffer{}
		for {
			cur, err := r.currentByte()
			if err != nil {
				return err
			}

			switch cur {
			case '\n':
				str.WriteByte('\n')
			default:
				str.WriteByte(cur)
			}

			peekChar, err := r.peek()

			if err != nil {
				if err.Inspect() == "EOF" {
					return &object.Error{Value: errors.New("Missing closing \"")}
				}

				return err
			}

			if peekChar == '"' {
				break
			}
		}

		cur, err := r.currentByte()

		if err != nil {
			return err
		}

		if cur != '"' {
			panic("Missing closing \"")
		}

		return &object.String{Value: str.String()}
	case '\'':
		cdr := r.Read()

		if isError(cdr) {
			return cdr
		}

		return &object.Cons{Car: &object.Identifier{Value: "quote"}, Cdr: &object.Cons{Car: cdr}}
	case '`':
		cdr := r.Read()

		if isError(cdr) {
			return cdr
		}

		return &object.Cons{Car: &object.Identifier{Value: "quasiquote"}, Cdr: &object.Cons{Car: cdr}}
	case '(':
		peekChar, err := r.peek()
		if err != nil {
			return err
		}

		if peekChar == ')' {
			return &object.Error{Value: errors.New("Unexpected closing )")}
		}

		obj := r.Read()
		if isError(obj) {
			return obj
		}

		list := &object.Cons{Car: obj}
		lastCons := list

		for {
			peekChar, err := r.peek()

			if err != nil {
				return err
			}

			if peekChar == ')' {
				break
			}

			cur, err := r.currentByte()
			if err != nil {
				return err
			}

			if cur == '.' {
				obj = r.Read()

				if isError(obj) {
					return obj
				}

				lastCons.Cdr = obj
			} else {
				err := r.unreadByte()
				if err != nil {
					return err
				}
				obj = r.Read()
				if isError(obj) {
					return obj
				}

				lastCons.Cdr = &object.Cons{Car: obj}
				lastCons = lastCons.Cdr.(*object.Cons)
			}
		}

		r.skip()

		return list
	case ' ', '\n', '\r', '\t':
		peekChar, err := r.peek()
		if err != nil {
			return err
		}

		// Consume white space
		for isWS(peekChar) {
			peekChar, err = r.peek()
			if err != nil {
				return err
			}
		}

		return r.Read()
	default:
		panic(fmt.Sprintf("No Clue: %#v", string(char)))
	}
}

func (r *Reader) peek() (byte, *object.Error) {
	bytes, err := r.reader.Peek(1)

	if err != nil {
		return byte(1), &object.Error{Value: err}
	}

	for isWS(bytes[0]) && err == nil {
		r.skip()
		return r.peek()
	}

	return bytes[0], nil
}

func (r *Reader) currentByte() (byte, *object.Error) {
	val, err := r.reader.ReadByte()

	if err != nil {
		return 1, &object.Error{Value: err}
	}

	return val, nil
}

func (r *Reader) unreadByte() *object.Error {
	err := r.reader.UnreadByte()

	if err != nil {
		return &object.Error{Value: err}
	}
	return nil
}

func (r *Reader) skip() {
	r.reader.Discard(1)
}

func isWS(char byte) bool {
	return ' ' == char || '\n' == char || '\r' == char || char == '\t'
}
