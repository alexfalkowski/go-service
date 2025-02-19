package metrics

const (
	// UnaryKind represents a unary RPC.
	UnaryKind Kind = "unary"

	// StreamKind represents a streaming RPC.
	StreamKind Kind = "stream"
)

// Kind represents the type of a metric.
type Kind string

// String representation of the kind.
func (k Kind) String() string {
	return string(k)
}
