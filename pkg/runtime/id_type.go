package runtime

type IDType interface {
	GetLiteral() string
}

type IDName struct {
	Literal string
}

func (id *IDName) GetLiteral() string {
	return id.Literal
}

type IDNumber struct {
	Literal  string
	NumValue float64
}

func (id *IDNumber) GetLiteral() string {
	return id.Literal
}

func (id *IDNumber) GetValue() float64 {
	return id.NumValue
}

func NewIDName(name string) *IDName {
	return &IDName{
		Literal: name,
	}
}
