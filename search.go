package trie

type Visit byte

const (
	Preorder Visit = iota
	Postorder
)

// TODO: might be nice to have a custom String() implementation here that includes the search metadata fields
type SearchNode struct {
	*Node
	Height      int
	ParentIndex int
	Order       Visit
}

// Traverse the sub-tree rooted at node in depth-first search order. The callback is called with each node twice: once
// when it is discovered (i.e. in preorder) and once when it is all of its children have been fully visited
// (i.e. in postorder). The callback can use SearchNode.Order to determine on which visit the callback is being invoked.
func DepthFirst(node *Node, callback func(SearchNode) error) error {
	stack := []*SearchNode{{Node: node}}
	for len(stack) > 0 {
		// pop
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if n.Order == Preorder {
			// visit in preorder
			err := callback(*n)
			if err != nil {
				return err
			}
			// set the next visit to be postorder
			n.Order = Postorder
			// push self
			stack = append(stack, n)
			// push children
			stack = push(stack, n)
		} else {
			err := callback(*n)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Traverse the sub-tree rooted at node in breadth-first search order. The callback is called with each node once in the
// order they are encountered - in monotonically increasing tree height heights.
func BreadthFirst(node *Node, callback func(SearchNode) error) error {
	queue := push(nil, &SearchNode{Node: node})
	for len(queue) > 0 {
		// pop
		n := queue[0]
		queue = queue[1:]
		// visit
		err := callback(*n)
		if err != nil {
			return err
		}
		// enqueue
		queue = push(queue, n)
	}
	return nil
}

// Push children of node on to dequeue slice
func push(dequeue []*SearchNode, node *SearchNode) []*SearchNode {
	height := node.Height + 1
	for _, child := range node.Children() {
		dequeue = append(dequeue, &SearchNode{Height: height, ParentIndex: node.index, Node: child})
	}
	return dequeue
}
