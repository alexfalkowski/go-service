package test

// AnyTuple for test.
type AnyTuple [2]any

// StringTuple for test.
type StringTuple [2]string

// KeyValue for test.
type KeyValue[Key, Value any] struct {
	Key   Key
	Value Value
}
