package trie

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: convert these into actual tests, probably using testing/quick
func TestDepthFirst(t *testing.T) {
	trie := NewTrie()
	set(trie, "hello")
	set(trie, "hellp")
	set(trie, "helicopter")
	set(trie, "hellpl")
	set(trie, "apple")
	set(trie, "arnold")
	set(trie, "butter")
	set(trie, "buttercup")
	fmt.Println(trie.Dump())

	fmt.Println("Preorder:")
	err := DepthFirst(trie.Node, func(node SearchNode) error {
		if node.Order == Preorder {
			fmt.Println(node)
		}
		return nil
	})
	require.NoError(t, err)

	fmt.Println()
	fmt.Println("Postorder:")
	err = DepthFirst(trie.Node, func(node SearchNode) error {
		if node.Order == Postorder {
			fmt.Println(node)
		}
		return nil
	})
	require.NoError(t, err)
}

func TestBreadthFirst(t *testing.T) {
	trie := NewTrie()
	set(trie, "hello")
	set(trie, "hellp")
	set(trie, "helicopter")
	set(trie, "hellpl")
	set(trie, "apple")
	set(trie, "arnold")
	set(trie, "butter")
	set(trie, "buttercup")
	fmt.Println(trie.Dump())

	err := BreadthFirst(trie.Node, func(node SearchNode) error {
		fmt.Println(node)
		return nil
	})
	require.NoError(t, err)
}
