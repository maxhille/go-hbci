package token

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	eof                       = -1
	dataElementDelimiter      = '+'
	groupDataElementDelimiter = ':'
	segmentDelimiter          = '\''
	escapeCharacter           = '?'
	binaryIdentifier          = '@'
)

// StringLexerStateFn represents a state function for the lexer.
type StringLexerStateFn func(*StringLexer) StringLexerStateFn

// NewStringLexer creates a new scanner for the input string.
func NewStringLexer(name, input string) *StringLexer {
	l := &StringLexer{
		name:   name,
		input:  input,
		state:  lexText,
		tokens: make(chan Token, 2), // Two token sufficient.
	}
	return l
}

// A StringLexer is a HBCI data element lexer based on an input string
type StringLexer struct {
	name   string             // the name of the input; used only for error reports.
	input  string             // the string being scanned.
	state  StringLexerStateFn // the next lexing function to enter
	pos    int                // current position in the input.
	start  int                // start position of this item.
	width  int                // width of last rune read from input.
	tokens chan Token         // channel of scanned tokens.
}

// SetEntryPoint sets the initial state of the lexer. The lexer will reset itself to use the
// new entryPoint properly
func (l *StringLexer) SetEntryPoint(entryPoint StringLexerStateFn) {
	l.reset()
	l.state = entryPoint
}

func (l *StringLexer) reset() {
	l.pos = 0
	l.start = 0
	l.width = 0
	l.tokens = make(chan Token, 2)
}

func (l *StringLexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens) // No more tokens will be delivered.
}

// Next returns the next item from the input.
func (l *StringLexer) Next() Token {
	for {
		select {
		case item, ok := <-l.tokens:
			if ok {
				return item
			}

			panic(fmt.Errorf("No items left"))
		default:
			l.state = l.state(l)
			if l.state == nil {
				close(l.tokens)
			}
		}
	}
}

// HasNext returns true if there are tokens left, false if EOF has reached
func (l *StringLexer) HasNext() bool {
	return l.state != nil
}

// emit passes a token back to the client.
func (l *StringLexer) emit(t Type) {
	l.tokens <- New(t, l.input[l.start:l.pos], l.start)
	l.start = l.pos
}

// next returns the next rune in the input.
func (l *StringLexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// ignore skips over the pending input before this point.
func (l *StringLexer) ignore() {
	l.start = l.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *StringLexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input.
func (l *StringLexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// accept consumes the next rune
// if it's from the valid set.
func (l *StringLexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *StringLexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on. Doing it this way
// means we don't have to worry about peek double counting.
func (l *StringLexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.pos], "\n")
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.run.
func (l *StringLexer) errorf(format string, args ...interface{}) StringLexerStateFn {
	l.tokens <- New(ERROR, fmt.Sprintf(format, args...), l.start)
	return nil
}

// state functions

func lexText(l *StringLexer) StringLexerStateFn {
	switch r := l.next(); {
	case r == dataElementDelimiter:
		l.emit(DATA_ELEMENT_SEPARATOR)
		return lexText
	case r == segmentDelimiter:
		l.emit(SEGMENT_END_MARKER)
		return lexText
	case r == groupDataElementDelimiter:
		l.emit(GROUP_DATA_ELEMENT_SEPARATOR)
		return lexText
	case r == binaryIdentifier:
		l.backup()
		return lexBinaryData
	case r == eof:
		// Correctly reached EOF.
		l.emit(EOF)
		return nil
	case ('0' <= r && r <= '9'):
		l.backup()
		return lexDigit
	default:
		l.backup()
		return lexAlphaNumeric
	}
}

func lexAlphaNumeric(l *StringLexer) StringLexerStateFn {
	text := false
	for {
		switch r := l.next(); {
		case r == escapeCharacter:
			if p := l.peek(); isSyntaxSymbol(p) {
				l.next()
			} else {
				return l.errorf("Unexpected escape character")
			}
		case isDelimiter(r):
			l.backup()
			if text {
				l.emit(TEXT)
			} else {
				l.emit(ALPHA_NUMERIC)
			}
			return lexText
		case r == eof:
			return l.errorf("Unexpected end of input")
		case (r == '\n' || r == '\r'):
			text = true
		}
	}
}

func lexBinaryData(l *StringLexer) StringLexerStateFn {
	l.accept("@")
	digits := "0123456789"
	binaryLengthStart := l.pos
	l.acceptRun(digits)
	binaryLengthEnd := l.pos
	if binaryLengthEnd == binaryLengthStart {
		return l.errorf("Binary length can't be empty")
	}
	length, err := strconv.Atoi(l.input[binaryLengthStart:binaryLengthEnd])
	if err != nil {
		return l.errorf("Binary length must contain of digits only")
	}
	if !l.accept("@") {
		return l.errorf("Binary length must contain of digits only")
	}
	l.pos += length
	p := l.peek()
	if isDelimiter(p) {
		l.emit(BINARY_DATA)
		return lexText
	}
	if p == eof {
		return l.errorf("Unexpected end of input")
	}
	return l.errorf("Expected syntax symbol after binary data")
}

func lexDigit(l *StringLexer) StringLexerStateFn {
	digits := "0123456789"
	leadingZero := l.accept("0")
	if leadingZero {
		// Only valid number with leading 0 is 0
		if r := l.peek(); isDelimiter(r) {
			l.emit(NUMERIC)
			return lexText
		}
		// Only valid float with leading 0 is value smaller than 1
		if l.accept(",") {
			l.acceptRun(digits)
			if p := l.peek(); isDelimiter(p) {
				l.emit(FLOAT)
				return lexText
			}
			return lexAlphaNumeric
		}
		l.acceptRun(digits)
		if p := l.peek(); p == ',' {
			return l.errorf("Malformed float")
		}
		if p := l.peek(); isDelimiter(p) {
			l.emit(DIGIT)
			return lexText
		}
		return lexAlphaNumeric
	}
	l.acceptRun(digits)
	// is it a float?
	if l.accept(",") {
		l.acceptRun(digits)
		if p := l.peek(); isDelimiter(p) {
			l.emit(FLOAT)
			return lexText
		}
		return lexAlphaNumeric
	}
	if p := l.peek(); isDelimiter(p) {
		l.emit(NUMERIC)
		return lexText
	}
	return lexAlphaNumeric

}

func isDelimiter(r rune) bool {
	return r == dataElementDelimiter || r == segmentDelimiter || r == groupDataElementDelimiter
}

func isSyntaxSymbol(r rune) bool {
	return r == dataElementDelimiter || r == segmentDelimiter || r == groupDataElementDelimiter || r == escapeCharacter || r == binaryIdentifier
}
