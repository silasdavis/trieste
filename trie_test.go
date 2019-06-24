package trie

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"runtime/debug"
	"sort"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	trie := NewTrie()
	set(trie, "hello")
	set(trie, "hellp")
	set(trie, "hellpl")
	set(trie, "apple")
	set(trie, "arnold")
	set(trie, "butter")
	set(trie, "buttercup")
	fmt.Println(trie.Dump())
}

func TestDelete(t *testing.T) {
	trie := NewTrie()
	set(trie, "hello")
	set(trie, "hellp")
	set(trie, "hellpl")
	set(trie, "apple")
	set(trie, "arnold")
	set(trie, "butter")
	set(trie, "buttercup")
	del(trie, "hellpl")
	del(trie, "hellp")
	set(trie, "hellp")
	fmt.Println(trie.Dump())
}

func TestDelete2(t *testing.T) {
	trie := NewTrie()
	set(trie, "hello")
	del(trie, "hello")
	set(trie, "hello")
	set(trie, "help")
	fmt.Println(trie.Dump())
}

func TestEmpty(t *testing.T) {
	trie := NewTrie()
	set(trie, "")
	fmt.Println(trie.Dump())
}

func TestGrow(t *testing.T) {
	tb := NewBranch(2)
	branchBytes := []byte("axcdpjks")
	for i := 0; i < len(branchBytes); i++ {
		bb := branchBytes[i]
		tb.Add(Char(bb), NewBranch(int(bb)))
	}
	require.Len(t, tb.children, len(branchBytes))
	sort.Slice(branchBytes, func(i, j int) bool {
		return branchBytes[i] < branchBytes[j]
	})
	for i := 0; i < len(branchBytes); i++ {
		index := tb.children[i].index
		require.Equal(t, int(branchBytes[i]), index)
	}
}

func TestPrefixInsert(t *testing.T) {
	trie := NewTrie()
	key := []byte("fos")
	trie.Set(key, key)
	key = []byte("f")
	trie.Set(key, key)
}

func TestDoubleInsert(t *testing.T) {
	trie := NewTrie()
	key := []byte{0}
	updated := trie.Set(key, key)
	assert.False(t, updated)
	updated = trie.Set(key, key)
	assert.True(t, updated)
}

func TestCriticalIndex(t *testing.T) {
	assert.Equal(t, 0, findCriticalIndex([]byte(""), []byte("")))
	assert.Equal(t, 0, findCriticalIndex([]byte("a"), []byte("b")))
	assert.Equal(t, 1, findCriticalIndex([]byte("aa"), []byte("ab")))
	assert.Equal(t, 2, findCriticalIndex([]byte("aa"), []byte("aa")))
	assert.Equal(t, 5, findCriticalIndex([]byte("aabra"), []byte("aabracadabra")))
	assert.Equal(t, 5, findCriticalIndex([]byte("aabracadabra"), []byte("aabra")))
	assert.Equal(t, 5, findCriticalIndex([]byte("aabra"), []byte("aabra")))
}

func TestChildIndex(t *testing.T) {
	tb := NewBranch(0)
	assert.Equal(t, uint(0), tb.childIndex('x'))
	tb.bitmap.Set('a')
	assert.Equal(t, uint(1), tb.childIndex('x'))
	tb.bitmap.Set('y')
	assert.Equal(t, uint(1), tb.childIndex('x'))
	tb.bitmap.Set('w')
	assert.Equal(t, uint(2), tb.childIndex('x'))
	tb.bitmap.Set('x')
	assert.Equal(t, uint(2), tb.childIndex('x'))
}

const maxKeyTapeLength = 1000000

type keytape struct {
	tape []byte
	rand *rand.Rand
}

func (kt *keytape) Keys() [][]byte {
	var keys [][]byte
	// cut the tape into random keys
	tape := make([]byte, len(kt.tape))
	copy(tape, kt.tape)
	for len(tape) > 0 {
		bite := kt.rand.Intn(len(tape))
		if bite == 0 {
			bite++
		}
		keys = append(keys, tape[:bite])
		tape = tape[bite:]
	}
	return keys
}

func newKeyTape(rand *rand.Rand, length int) *keytape {
	tape := make([]byte, length)
	for i := 0; i < length; i++ {
		tape[i] = byte(rand.Intn(256))
	}
	return &keytape{
		rand: rand,
		tape: tape,
	}
}

func keytapeValues(args []reflect.Value, rand *rand.Rand) {
	for i := range args {
		args[i] = reflect.ValueOf(newKeyTape(rand, rand.Intn(maxKeyTapeLength)+1))
	}
}

func TestKeyTape(t *testing.T) {
	kt := newKeyTape(rand.New(rand.NewSource(32423)), 100)
	keys := kt.Keys()
	fmt.Println(keys)
}

func TestConsistency(t *testing.T) {
	var err error
	var trie *Trie
	var badKey []byte
	setThenGet := func(kt *keytape) (success bool) {
		defer func() {
			if r := recover(); r != nil {
				success = false
				err = fmt.Errorf("%v:\n%s", r, string(debug.Stack()))
			}
		}()
		trie = NewTrie()
		keys := kt.Keys()
		for _, key := range keys {
			badKey = key
			trie.Set(key, key)
		}
		for _, key := range keys {
			badKey = key
			value, exists := trie.Get(key)
			if !exists {
				err = fmt.Errorf("key %X was inserted but could not be retrieved", key)
				return false
			}
			if !bytes.Equal(key, value.([]byte)) {
				err = fmt.Errorf("inserted key/value %X, but retrieved %X", key, value.([]byte))
				return false
			}
		}
		return true
	}
	success := setThenGet(&keytape{
		tape: []byte{},
		rand: rand.New(rand.NewSource(242323)),
	})

	require.True(t, success)
	checkErr := quick.Check(setThenGet, &quick.Config{
		Values: keytapeValues,
	})
	if checkErr != nil {
		//t.Log(trie.Dump())
		t.Logf("Bad key: %X", badKey)
		trie.Set(badKey, badKey)
		t.Error(err)
	}
}

func BenchmarkArrayAccess(b *testing.B) {
	key := []byte("alongishkindofkey")
	keycopy := make([]uint, len(key))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(key); j++ {
			keycopy[j] = uint(key[j])
		}
	}
}

func BenchmarkCharAt(b *testing.B) {
	key := []byte("alongishkindofkey")
	keycopy := make([]Char, len(key)+1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 1; j <= len(key); j++ {
			keycopy[j] = charAt(key, j)
		}
	}
}

func set(trie *Trie, key string) bool {
	return trie.Set([]byte(key), key)
}

func del(trie *Trie, key string) bool {
	return trie.Delete([]byte(key))
}
