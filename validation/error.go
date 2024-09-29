package validation

import (
	"strconv"
	"strings"
)

type RuleError interface {
	Error() string
	Code() string
	Message() string
}

type ruleError struct {
	code    string
	message string
}

func NewRuleError(code, message string) RuleError {
	return &ruleError{
		code:    code,
		message: message,
	}
}

func (re *ruleError) Error() string {
	return re.message
}

func (re *ruleError) Code() string {
	return re.code
}

func (re *ruleError) Message() string {
	return re.message
}

type ValueError interface {
	Error() string
	Unwrap() error
	Name() string
}

type valueError struct {
	name   string
	nested error
}

func NewValueError(name string, nested error) ValueError {
	return &valueError{
		name:   name,
		nested: nested,
	}
}

func (ve *valueError) Error() string {
	msg := ve.nested.Error()

	var b strings.Builder

	b.Grow(len(ve.name) + 2 + len(msg))
	b.WriteString(ve.name)
	b.WriteString(": ")
	b.WriteString(msg)

	return b.String()
}

func (ve *valueError) Name() string {
	return ve.name
}

func (ve *valueError) Unwrap() error {
	return ve.nested
}

type IndexError interface {
	Error() string
	Unwrap() error
	Index() int
}

type indexError struct {
	index  int
	nested error
}

func NewIndexError(index int, nested error) IndexError {
	return &indexError{
		index:  index,
		nested: nested,
	}
}

func (ie *indexError) Error() string {
	msg := ie.nested.Error()
	idx := strconv.Itoa(ie.index)

	var b strings.Builder

	b.Grow(1 + len(idx) + 3 + len(msg))
	b.WriteByte('[')
	b.WriteString(idx)
	b.WriteString("]: ")
	b.WriteString(msg)

	return b.String()
}

func (ie *indexError) Index() int {
	return ie.index
}

func (ie *indexError) Unwrap() error {
	return ie.nested
}

type Errors []error

func errorWriteString(err error, b *strings.Builder) {
	switch e := err.(type) {
	case Errors:
		b.WriteByte('(')
		e.writeString(b)
		b.WriteByte(')')
	case IndexError:
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(e.Index()))
		b.WriteString(": ")
		errorWriteString(e.Unwrap(), b)
		b.WriteByte(')')
	default:
		b.WriteString(e.Error())
	}
}

func (es Errors) writeString(b *strings.Builder) {
	sz := b.Len()
	for _, err := range es {
		switch e := err.(type) {
		case Errors:
			e.writeString(b)
		case ValueError:
			if sz != b.Len() {
				b.WriteString("; ")
			}
			b.WriteString(e.Name())
			b.WriteString(": ")
			errorWriteString(e.Unwrap(), b)
		}
	}
}

func (es Errors) Error() string {
	if len(es) == 0 {
		return ""
	}
	var b strings.Builder
	es.writeString(&b)
	return b.String()
}

func errorMarshalJSON(err error, b []byte) []byte {
	switch e := err.(type) {
	case Errors:
		b = append(b, '{')
		b = e.marshalJSON(b)
		b = append(b, '}')
	case IndexError:
		b = append(b, '{', '"')
		b = strconv.AppendInt(b, int64(e.Index()), 10)
		b = append(b, '"', ':')
		b = errorMarshalJSON(e.Unwrap(), b)
		b = append(b, '}')
	default:
		b = strconv.AppendQuote(b, e.Error())
	}
	return b
}

func (es Errors) marshalJSON(b []byte) []byte {
	sz := len(b)
	for _, err := range es {
		switch e := err.(type) {
		case Errors:
			b = e.marshalJSON(b)
		case ValueError:
			if sz != len(b) {
				b = append(b, ',')
			}
			b = strconv.AppendQuote(b, e.Name())
			b = append(b, ':')
			b = errorMarshalJSON(e.Unwrap(), b)
		}
	}
	return b
}

func (es Errors) MarshalJSON() ([]byte, error) {
	var b []byte
	b = append(b, '{')
	b = es.marshalJSON(b)
	b = append(b, '}')
	return b, nil
}
