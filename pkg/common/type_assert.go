package common

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

// ////// Either Type (2 Types) ///////
type Either[A, B any] struct {
	value any
}

// before use GetA(), make sure IsA() is true!!!
func (e *Either[A, B]) GetA() A {
	return e.value.(A)
}

// before use GetB(), make sure IsB() is true!!!
func (e *Either[A, B]) GetB() B {
	return e.value.(B)
}

func (e *Either[A, B]) IsA() bool {
	_, ok := e.value.(A)
	return ok
}

func (e *Either[A, B]) IsB() bool {
	_, ok := e.value.(B)
	return ok
}

// ////// Either3 Type (3 Types) ///////
type Either3[A, B, C any] struct {
	value any
}

// before use GetA(), make sure IsA() is true!!!
func (e *Either3[A, B, C]) GetA() A {
	return e.value.(A)
}

// before use GetB(), make sure IsB() is true!!!
func (e *Either3[A, B, C]) GetB() B {
	return e.value.(B)
}

func (e *Either3[A, B, C]) GetC() C {
	return e.value.(C)
}

func (e *Either3[A, B, C]) IsA() bool {
	_, ok := e.value.(A)
	return ok
}

func (e *Either3[A, B, C]) IsB() bool {
	_, ok := e.value.(B)
	return ok
}

func (e *Either3[A, B, C]) IsC() bool {
	_, ok := e.value.(C)
	return ok
}

// /// ASSERT ELEMENT TYPE /////
func AssertElement[T r.Element](elem r.Element) (T, error) {
	if v, ok := elem.(T); ok {
		return v, nil
	}
	return elem.(T), zerr.InvalidParamType(getElementTypeString(elem))
}

func AssertInstanceOf(elem r.Element, model *value.ClassModel) (*value.Object, error) {
	obj, ok := elem.(*value.Object)
	if !ok {
		return nil, zerr.InvalidParamType(getElementTypeString(elem))
	}
	if !obj.IsInstanceOf(model) {
		return nil, zerr.InvalidParamType(model.String())
	}
	return obj, nil

}

func AssertPropertyElement[T r.Element](root r.Element, key string) (T, error) {
	prop, err := root.GetProperty(key)
	if err != nil {
		return root.(T), err
	}

	if v, ok := prop.(T); ok {
		return v, nil
	}
	return root.(T), zerr.InvalidParamType(getElementTypeString(prop))
}

func AssertEitherElement[A r.Element, B r.Element](elem r.Element) (Either[A, B], error) {
	if _, ok := elem.(A); !ok {
		if _, ok := elem.(B); !ok {
			return Either[A, B]{nil}, zerr.InvalidParamType(getElementTypeString(elem))
		}
	}
	return Either[A, B]{elem}, nil
}

func AssertEitherPropertyElement[A r.Element, B r.Element](root r.Element, key string) (Either[A, B], error) {
	prop, err := root.GetProperty(key)
	if err != nil {
		return Either[A, B]{nil}, err
	}
	if _, ok := prop.(A); !ok {
		if _, ok := prop.(B); !ok {
			return Either[A, B]{nil}, zerr.InvalidParamType(getElementTypeString(prop))
		}
	}
	return Either[A, B]{prop}, nil
}

func getElementTypeString(elem r.Element) string {
	switch elem.(type) {
	case *value.Number:
		return "number"
	case *value.String:
		return "string"
	case *value.Array:
		return "array"
	case *value.HashMap:
		return "hashmap"
	case *value.Bool:
		return "bool"
	case *value.Object:
		return "object"
	case *value.Function:
		return "function"
	case *value.GoValue:
		return "govalue"
	}
	return "unknown"
}
