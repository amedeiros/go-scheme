package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// TRUE is the only true value
var TRUE = &Boolean{Value: true}

// FALSE is the only false value
var FALSE = &Boolean{Value: false}

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

// ReadAll will read until it encounters an error
func (r *Reader) ReadAll() []Object {
	var program []Object

	for {
		obj := r.Read()

		if isError(obj) {
			err := obj.(*Error)
			if err.Value.Error() != EOF {
				program = append(program, err)
			}

			break
		}

		program = append(program, obj)
	}

	return program
}

// Read will parse and return an object on each call
func (r *Reader) Read() Object {
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
			return &Char{Value: string(cur)}
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

			peekChar, err := r.preserveWsPeek(true)

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

		return &String{Value: str.String()}
	case '\'':
		cdr := r.Read()

		if isError(cdr) {
			return cdr
		}

		return &Cons{Car: &Identifier{Value: "QUOTE"}, Cdr: &Cons{Car: &String{Value: cdr.Inspect()}}}
	case '`':
		cdr := r.Read()

		if isError(cdr) {
			return cdr
		}

		return &Cons{Car: &Identifier{Value: "quasiquote"}, Cdr: &Cons{Car: cdr}}
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

		// if obj.Type() == IDENT_OBJ {
		// 	ident := obj.(*Identifier)
		// 	if ident.Value == "LAMBDA" {
		// 		return r.readLambda()
		// 	}
		// }

		switch node := obj.(type) {
		case *Identifier:
			if node.Value == "LAMBDA" {
				return r.readLambda()
			}
		}

		list := &Cons{Car: obj}
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

				lastCons.Cdr = &Cons{Car: obj}
				lastCons = lastCons.Cdr.(*Cons)
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
	case '+', '*', '/', '=':
		return &Identifier{Value: string(char)}
	case '-':
		return r.identOrDigit(char)
	case ';':
		r.consumeComment()
		return r.Read()
	default:
		return r.identOrDigit(char)
	}
}

func (r *Reader) consumeComment() {
	peekChar, _ := r.preserveWsPeek(true)

	for peekChar != '\n' && peekChar != '\r' {
		r.skip()
		peekChar, _ = r.preserveWsPeek(true)
	}
}

func (r *Reader) identOrDigit(char byte) Object {
	str := bytes.Buffer{}
	str.WriteByte(char)
	for {
		char, err := r.currentByte()
		if err != nil {
			if err.Value.Error() == EOF {
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
			return &Identifier{Value: strings.ToUpper(str.String())}
		}

		return &Float{Value: f}
	}

	return &Integer{Value: i}
}

func (r *Reader) readLambda() Object {
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

	var body Object

	if curChar == ')' {
		// No Arguments
		body = r.Read()

		curChar, _ = r.currentByte()

		if curChar != ')' {
			return newError("missing closing )")
		}

		return &Lambda{Body: body}
	}

	var arguments []*Identifier

	r.unreadByte() // Go back one or we skip the first identifier

	for {
		peekChar, err := r.peek()
		if err != nil {
			return err
		}

		if peekChar == ')' {
			break
		}

		arg := r.Read()
		arguments = append(arguments, arg.(*Identifier))
	}

	r.skip() // Skip closing )
	body = r.Read()
	r.skip() // Skip closing )

	if isError(body) {
		return body
	}

	return &Lambda{Body: body, Parameters: arguments}
}

func (r *Reader) peek() (byte, *Error) {
	return r.preserveWsPeek(false)
}

func (r *Reader) preserveWsPeek(enable bool) (byte, *Error) {
	bytes, err := r.reader.Peek(1)

	if err != nil {
		return 1, &Error{Value: err}
	}

	if !enable && isWS(bytes[0]) {
		for isWS(bytes[0]) {
			r.skip()
			bytes, err = r.reader.Peek(1)
			if err != nil {
				return 1, errorObject(err)
			}
		}
	}

	return bytes[0], nil
}

func (r *Reader) currentByte() (byte, *Error) {
	val, err := r.reader.ReadByte()

	if err != nil {
		return 1, &Error{Value: err}
	}

	return val, nil
}

func (r *Reader) unreadByte() *Error {
	err := r.reader.UnreadByte()

	if err != nil {
		return &Error{Value: err}
	}
	return nil
}

func (r *Reader) skip() {
	r.reader.Discard(1)
}

func isWS(char byte) bool {
	return ' ' == char || '\n' == char || '\r' == char || char == '\t'
}
