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
		} else if peekChar == "(" {
			r.skip()
			values := []Object{}

			for {
				peekChar, err := r.peek()
				if err != nil {
					return err
				}

				if peekChar == ')' {
					r.skip()
					break
				}

				values = append(values, r.Read())
			}

			return &Vector{Value: values}
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
			case '"':
				return &String{Value: ""}
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

		return &Pair{Car: &Identifier{Value: "QUOTE"}, Cdr: cdr}
	case '`':
		cdr := r.Read()

		if isError(cdr) {
			return cdr
		}

		return &Pair{Car: &Identifier{Value: "QUASIQUOTE"}, Cdr: &Pair{Car: cdr}}
	case '(':
		peekChar, err := r.peek()

		if err != nil {
			return err
		}

		if peekChar == ')' {
			r.skip()
			return &Pair{}
		}

		obj := r.Read()
		if isError(obj) {
			return obj
		}

		switch node := obj.(type) {
		case *Identifier:
			if node.Value == "LAMBDA" {
				return r.readLambda()
			} else if node.Value == "LET" {
				return r.expandLet()
			} else if node.Value == "DEFINE" {
				return r.expandDefine(node)
			}
		}

		list := &Pair{Car: obj}
		lastPair := list

		for {
			peekChar, err := r.peek()

			if err != nil {
				return err
			}

			if peekChar == ')' {
				r.skip()
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

				lastPair.Cdr = obj
			} else {
				err := r.unreadByte()
				if err != nil {
					return err
				}
				obj = r.Read()
				if isError(obj) {
					return obj
				}

				lastPair.Cdr = &Pair{Car: obj}
				lastPair = lastPair.Cdr.(*Pair)
			}
		}

		return list
	case ' ', '\n', '\r', '\t':
		peekChar, err := r.peek()
		if err != nil {
			return err
		}

		// Pairume white space
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
	case '<':
		peekChar, err := r.peek()
		if err != nil {
			return err
		}

		if peekChar == '=' {
			r.skip()
			return &Identifier{Value: "<="}
		}

		return &Identifier{Value: "<"}
	case '>':
		peekChar, err := r.peek()
		if err != nil {
			return err
		}

		if peekChar == '=' {
			r.skip()
			return &Identifier{Value: ">="}
		}

		return &Identifier{Value: ">"}
	case ';':
		r.parseComment()
		return r.Read()
	default:
		return r.identOrDigit(char)
	}
}

// Expand define into something out interpreter can handle
func (r *Reader) expandDefine(ident *Identifier) Object {
	pair := &Pair{Car: ident}

	first := r.Read()
	switch kind := first.(type) {
	case *Identifier:
		value := r.Read()
		pair.Cdr = &Pair{Car: kind, Cdr: &Pair{Car: value}}
	case *Pair:
		variable := car(kind)
		params := []*Identifier{}

		for {
			param := kind.Cdr.(*Pair)
			params = append(params, car(param).(*Identifier))

			if param.Cdr != nil {
				kind.Cdr = param.Cdr
			} else {
				break
			}

		}

		body := r.Read()
		lambda := &Pair{Car: &Lambda{Parameters: params, Body: body}}
		pair.Cdr = &Pair{Car: variable, Cdr: lambda}
	}

	return pair
}

// Let expands into a lambda call
func (r *Reader) expandLet() Object {
	args := []string{}
	params := []string{}

	if pair, ok := r.Read().(*Pair); ok {
		for {
			if first, ok := car(pair).(*Pair); ok {
				if ident, ok := car(first).(*Identifier); ok {
					params = append(params, ident.Inspect())
					arg := car(cdr(first))
					args = append(args, arg.Inspect())
				} else {
					pair = first
					continue
				}

				if pair.Cdr != nil {
					pair.Car = pair.Cdr
				} else {
					break
				}
			} else {
				return newError("expecting a proper list")
			}
		}
	} else {
		return newError("expecting a proper list")
	}

	body := r.Read()
	peekChar, err := r.peek()
	if err != nil {
		return err
	}

	if peekChar == ')' {
		r.skip()
	}

	// Expand let into a lambda call to preserve lexical scoping
	lambda := fmt.Sprintf("((lambda (%s) %s) %s)", strings.Join(params, " "), body.Inspect(), strings.Join(args, " "))
	reader := NewReader(lambda)
	return reader.Read()
}

func car(obj Object) Object {
	if pair, ok := obj.(*Pair); ok {
		return pair.Car
	}

	return newError("expecting a proper list")
}

func cdr(obj Object) Object {
	if pair, ok := obj.(*Pair); ok {
		return pair.Cdr
	}

	return newError("expecting a proper list")
}

func (r *Reader) parseComment() {
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
