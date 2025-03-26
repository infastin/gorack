package opt

import "encoding"

type Text interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

type ComparableText interface {
	comparable
	Text
}

type NullText[T Text] struct {
	Null[T]
}

func NewNullText[T Text](value T, valid bool) NullText[T] {
	return NullText[T]{
		Null: NewNull(value, valid),
	}
}

func NullTextFrom[T Text](value T) NullText[T] {
	return NullText[T]{
		Null: NullFrom(value),
	}
}

func NullTextFromPtr[T Text](value *T) NullText[T] {
	return NullText[T]{
		Null: NullFromPtr(value),
	}
}

func NullTextFromFunc[T Text, U any](value U, f func(U) T) NullText[T] {
	return NullText[T]{
		Null: NullFromFunc(value, f),
	}
}

func NullTextFromPtrFunc[T Text, U any](value *U, f func(U) T) NullText[T] {
	return NullText[T]{
		Null: NullFromPtrFunc(value, f),
	}
}

func (v NullText[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return v.Value.MarshalText()
}

func (v *NullText[T]) UnmarshalText(data []byte) error {
	if err := v.Value.UnmarshalText(data); err != nil {
		return err
	}
	v.Valid = true
	return nil
}

type ZeroText[T ComparableText] struct {
	Zero[T]
}

func NewZeroText[T ComparableText](value T, valid bool) ZeroText[T] {
	return ZeroText[T]{
		Zero: NewZero(value, valid),
	}
}

func ZeroTextFrom[T ComparableText](value T) ZeroText[T] {
	return ZeroText[T]{
		Zero: ZeroFrom(value),
	}
}

func ZeroTextFromPtr[T ComparableText](value *T) ZeroText[T] {
	return ZeroText[T]{
		Zero: ZeroFromPtr(value),
	}
}

func ZeroTextFromFunc[T ComparableText, U any](value U, f func(U) T) ZeroText[T] {
	return ZeroText[T]{
		Zero: ZeroFromFunc(value, f),
	}
}

func ZeroTextFromPtrFunc[T ComparableText, U any](value *U, f func(U) T) ZeroText[T] {
	return ZeroText[T]{
		Zero: ZeroFromPtrFunc(value, f),
	}
}

func (v ZeroText[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return v.Value.MarshalText()
}

func (v *ZeroText[T]) UnmarshalText(data []byte) error {
	if err := v.Value.UnmarshalText(data); err != nil {
		return err
	}

	var zero T
	v.Valid = v.Value != zero

	return nil
}
