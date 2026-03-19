package test

// AnyTuple is a generic two-slot tuple for loosely typed test fixtures.
type AnyTuple [2]any

// StringTuple is a two-slot tuple specialized for string values.
type StringTuple [2]string

// KeyValue stores a strongly typed key and value pair for helper assertions and fixtures.
type KeyValue[Key, Value any] struct {
	Key   Key
	Value Value
}
