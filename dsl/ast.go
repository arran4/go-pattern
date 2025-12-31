package dsl

import "fmt"

type Node interface {
	String() string
}

type CommandNode struct {
	Name string
	Args []ArgNode
}

func (n *CommandNode) String() string {
	s := n.Name
	for _, arg := range n.Args {
		s += " " + arg.String()
	}
	return s
}

type ArgNode interface {
	Node
	isArg()
}

type LiteralNode struct {
	Value string
}

func (n *LiteralNode) String() string {
	return n.Value
}
func (n *LiteralNode) isArg() {}

type KeyValueNode struct {
	Key   string
	Value ArgNode
}

func (n *KeyValueNode) String() string {
	return fmt.Sprintf("%s=%s", n.Key, n.Value.String())
}
func (n *KeyValueNode) isArg() {}

type PipelineNode struct {
	Nodes []Node
}

func (n *PipelineNode) String() string {
	s := ""
	for i, node := range n.Nodes {
		if i > 0 {
			s += " | "
		}
		s += node.String()
	}
	return s
}

type BinaryNode struct {
	Left     Node
	Operator string
	Right    Node
}

func (n *BinaryNode) String() string {
	return fmt.Sprintf("(%s) %s (%s)", n.Left.String(), n.Operator, n.Right.String())
}

type GroupNode struct {
	Inner Node
}

func (n *GroupNode) String() string {
	return fmt.Sprintf("(%s)", n.Inner.String())
}

// SubExpressionNode is used when an argument is a parenthesized expression, e.g. mask=(...)
type SubExpressionNode struct {
	Node Node
}

func (n *SubExpressionNode) String() string {
	return fmt.Sprintf("(%s)", n.Node.String())
}
func (n *SubExpressionNode) isArg() {}
