package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/amedeiros/go-scheme/object"
)

// TRUE is the only true value
var TRUE = &object.Boolean{Value: true}

// FALSE is the only false value
var FALSE = &object.Boolean{Value: false}

// EOF check for end of file
const EOF = "EOF"

// Reader wraps a bufio.Reader for us
type Reader struct {
	reader *bufio.Reader
}

// NewReader takes in a string and returns a new Reader
func NewReader(input string) *Reader {
	return &Reader{bufio.NewReader(strings.NewReader(input))}
}

// Read will parse and return an object on each call
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

		return newError(fmt.Sprintf("Expecting one of F or T or \\ found %s instead.", peekChar))
	case '"':
		str := bytes.Buffer{}

		for {
			cur, err := r.currentByte()

			if err != nil {
				return err
			}

			switch cur {
			case '\n':
				str.WriteString("\n")
			default:
				str.WriteByte(cur)
			}

			peekChar, err := r.peek()

			if err != nil {
				if err.Inspect() == EOF {
					return newError("Missing closing \"")
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
			return newError("Missing closing \"")
		}

		return &object.String{Value: str.String()}
	case '\'':
		cdr := r.Read()

		if isError(cdr) {
			return cdr
		}

		return &object.Cons{Car: &object.Identifier{Value: "QUOTE"}, Cdr: &object.Cons{Car: &object.String{Value: cdr.Inspect()}}}
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
			r.skip()
			return r.Read()
		}

		obj := r.Read()
		if isError(obj) {
			return obj
		}

		if obj.Type() == object.IDENT_OBJ {
			ident := obj.(*object.Identifier)
			if ident.Value == "LAMBDA" {
				return r.readLambda()
			}
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
	case '+', '*', '/', '-':
		return &object.Identifier{Value: string(char)}
	default:
		str := bytes.Buffer{}
		str.WriteByte(char)
		for {
			char, err = r.currentByte()
			if err != nil {
				if err.Value.Error() == "EOF" {
					break
				}

				return err
			}

			if isWS(char) {
				break
			}

			if char == '(' || char == ')' {
				err := r.unreadByte()
				if err != nil {
					return err
				}

				break
			}

			str.WriteByte(char)

		}

		i, err := strconv.ParseInt(str.String(), 0, 64)

		if err != nil {
			f, err := strconv.ParseFloat(str.String(), 64)

			if err != nil {
				return &object.Identifier{Value: strings.ToUpper(str.String())}
			}

			return &object.Float{Value: f}
		}

		return &object.Integer{Value: i}
	}
}

func (r *Reader) readLambda() object.Object {
	peekChar, err := r.peek()
	if err != nil {
		return err
	}

	if peekChar != '(' {
		return newError("Missing opening (")
	}

	r.skip()
	curChar, err := r.currentByte()
	if err != nil {
		return err
	}

	var body object.Object

	if curChar == ')' {
		// No Arguments
		body = r.Read()
		return &object.Lambda{Body: body}
	}

	err = r.unreadByte()
	if err != nil {
		return err
	}

	var arguments []*object.Identifier

	for {
		peekChar, err := r.peek()
		if err != nil {
			return err
		}

		if peekChar == ')' {
			r.skip()
			break
		}

		arg := r.Read()
		arguments = append(arguments, arg.(*object.Identifier))
	}

	peekChar, _ = r.peek()

	if peekChar == ')' {
		r.skip()
	}

	body = r.Read()

	if isError(body) {
		return body
	}

	return &object.Lambda{Body: body, Parameters: arguments}
}

func (r *Reader) peek() (byte, *object.Error) {
	bytes, err := r.reader.Peek(1)

	if err != nil {
		return byte(1), &object.Error{Value: err}
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
