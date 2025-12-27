package object

import "fmt"

type ObjectType string

const (
	IntegerType ObjectType = "INTEGER"
	BooleanType ObjectType = "BOOLEAN"
	NillType    ObjectType = "NILL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (integer Integer) Type() ObjectType {
	return IntegerType
}
func (integer Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

type Boolean struct {
	Value bool
}

func (boolean Boolean) Type() ObjectType {
	return BooleanType
}
func (boolean Boolean) Inspect() string {
	return fmt.Sprintf("%t", boolean.Value)
}

type Nill struct {
}

func (nill Nill) Type() ObjectType {
	return NillType
}
func (nill Nill) Inspect() string {
	return "null"
}
