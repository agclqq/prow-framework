package dag

type Builder interface {
	InsertBefore(node *Node) error
	InsertAfter(node *Node) error
	InsertBetween(node1, node2 *Node) error
}
type Verifier interface {
	Verify() error
}
type Driver interface {
	Drive() ([]*Node, bool)
}
