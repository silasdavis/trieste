package trie

import (
	"fmt"

	bitset "github.com/tmthrgd/go-bitset"
)

// Branch on byte plus special terminal character
const BranchingFactor uint = 257
const TerminalChar = 0

// Like a byte, but these ones go up to 256 (to make room for a distinguished terminal - 0
type Char uint

func CharFromByte(b byte) Char {
	return Char(b) + 1
}

// We model Node as a sum type. Either Node should be nil, contain just a non-nil Branch, or contain just a non-nil Leaf
type Node struct {
	*Branch
	*Leaf
}

// A Branch or internal node consists of an index which is the index of the key byte on which this branch 'decides'.
// A Branch implements a sparse array of children with one possible slot for each possible character branch.
// The Leaf node, Leaf[K,V], for a particular key can be found under the subtree rooted at the child in position C when
// charAt(K[Branch.index]) == C.
type Branch struct {
	// The critical index into keys on which the children of this branch are distinguished
	index int
	// The pointers to the children are over Chars' which are almost the key bytes but shifted by one to make space
	// for a distinguished terminal 0 (to provide a slot for empty keys and leaves that are a prefix of some branch)
	bitmap bitset.Bitset
	// TBC: chunking children into multiple arrays - e.g. 4 quadrants - to limit the effect of having to copy all the tail
	// of child pointers on insert
	children []*Node
}

type Leaf struct {
	// We need to store the entire key in a Leaf since the tree structure does not generally encode it fully
	// (just enough for an unambiguous match)
	key   []byte
	value interface{}
}

func NewBranch(index int) *Node {
	return &Node{
		Branch: &Branch{
			index:  index,
			bitmap: bitset.New(BranchingFactor),
		},
	}
}

func (tn *Node) IsLeaf() bool {
	return tn != nil && tn.Leaf != nil
}

func (tn *Node) IsBranch() bool {
	return tn != nil && tn.Branch != nil
}

func (tn *Node) String() string {
	if tn == nil {
		return "Node<nil>"
	}
	if tn.IsLeaf() {
		return fmt.Sprintf("Leaf<%X -> %v>", tn.key, tn.value)
	}
	return fmt.Sprintf("Branch<@ %d => [%d]>", tn.index, len(tn.children))
}

func (tn *Node) Children() []*Node {
	if tn.IsBranch() {
		return tn.children
	}
	return nil
}

// Descend implements the basic traversal of the trie. It returns the deepest directed edge,
// parent --char--> child, accessible by following key on each critical index along the path of Branches.
// child may be nil indicating the slot for char on branch is unoccupied.
// Note: we use the Node container for branch as a convenience for mutation - it is guaranteed to be a Branch Node
func (tn *Node) Descend(key []byte) (branch *Node, char Char, child *Node) {
	branch = tn
	char = charAt(key, branch.index)
	child = branch.lookup(char)
	for child.IsBranch() && child.index <= len(key) {
		branch = child
		char = charAt(key, branch.index)
		child = branch.lookup(char)
	}
	return
}

func (tb *Branch) Add(char Char, node *Node) {
	if node == nil {
		panic("tried to add a nil node")
	}
	childIndex := tb.childIndex(char)
	// insert node in between head and tail of children
	tb.children = append(tb.children[:childIndex], append([]*Node{node}, tb.children[childIndex:]...)...)
	tb.bitmap.Set(uint(char))
}

func (tb *Branch) Remove(char Char) {
	childIndex := tb.childIndex(char)
	// overwrite tail
	copy(tb.children[childIndex:], tb.children[childIndex+1:])
	// zero for safety
	tb.children[len(tb.children)-1] = nil
	tb.children = tb.children[:len(tb.children)-1]
	tb.bitmap.Clear(uint(char))
}

func (tb *Branch) lookup(char Char) *Node {
	if tb == nil {
		return nil
	}
	if tb.bitmap.IsSet(uint(char)) {
		return tb.children[tb.childIndex(char)]
	}
	return nil
}

func (tb *Branch) childIndex(char Char) uint {
	return tb.bitmap.CountRange(0, uint(char))
}

// Map key bytes to uint 'char' which includes a distinguished terminal 0
func charAt(key []byte, index int) Char {
	if 0 <= index && index < len(key) {
		return CharFromByte(key[index])
	}
	return TerminalChar
}

func byteFromChar(char Char) byte {
	return byte(char - 1)
}

func stringFromChars(chars ...Char) string {
	bs := make([]byte, len(chars))
	for i, char := range chars {
		bs[i] = byteFromChar(char)
	}
	return string(bs)
}

func (tb *Branch) childChars() []Char {
	var char Char
	chars := make([]Char, len(tb.children))
	for i := range chars {
		for !tb.bitmap.IsSet(uint(char)) {
			char++
		}
		chars[i] = char
		char++
	}
	return chars
}
