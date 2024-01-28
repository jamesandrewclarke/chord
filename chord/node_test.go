package chord

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingletonNodeIsOwnSuccessor(t *testing.T) {
	node := CreateNode(1)

	assert.Equal(t, node.Identifier(), node.successor.Identifier(), "A new node should be its own successor")
}

func TestFindSuccessorSimple(t *testing.T) {
	a := CreateNode(1)
	b := CreateNode(10)

	a.Join(b)

	for i := a.Identifier() + 1; i <= b.Identifier(); i++ {
		succ, err := a.FindSuccessor(Id(i))

		assert.Nil(t, err)
		assert.Equal(t, b.Identifier(), succ.Identifier())
	}
}

func TestFindSuccessorWrapAround(t *testing.T) {
	a := CreateNode(10)
	b := CreateNode(1)

	a.Join(b)

	succ, _ := a.FindSuccessor(a.Identifier())
	assert.Equal(t, a.Identifier(), succ.Identifier())

	// Any key >10 should be handled by node b
	for i := 1; i < 16; i++ {
		succ, _ = a.FindSuccessor(a.Identifier() + 1)
		assert.Equal(t, b.Identifier(), succ.Identifier())
	}
}

func TestSingletonNodeFindSuccessorReturnsSelf(t *testing.T) {
	node := CreateNode(1)

	for i := 1; i < 100; i++ {
		succ, err := node.FindSuccessor(node.Identifier())

		assert.Nil(t, err)
		assert.Equal(t, node.Identifier(), succ.Identifier())
	}
}

func TestJoinSetsCorrectSuccessor(t *testing.T) {
	a := CreateNode(1)
	b := CreateNode(10)

	a.Join(b)

	assert.Equal(t, b.Identifier(), a.successor.Identifier())

	// Stabilisation has not yet run, so b should not have a predecessor
	assert.Nil(t, b.predecessor)
}

func TestNotifySetsPredecessor(t *testing.T) {
	a := CreateNode(1)
	b := CreateNode(2)
	a.Join(b)

	b.Notify(a)

	assert.Equal(t, a.Identifier(), b.predecessor.Identifier(), "Predecessor should be set")
}

func TestNotifyRejectsInvalidPredecessor(t *testing.T) {
	a := CreateNode(2)
	b := CreateNode(4)

	a.Join(b)
	b.Notify(a)

	// Create a new node which comes before A. If B is notified by C, B's predecessor should still be A
	c := CreateNode(1)
	c.Join(a)
	b.Notify(c)

	assert.Equal(t, a.Identifier(), b.predecessor.Identifier(), "The current predecessor should be unchanged")
}

func TestStabilizeSetsSuccessor(t *testing.T) {
	a := CreateNode(1)
	b := CreateNode(2)
	a.Join(b)

	a.stabilize()

	assert.Equal(t, a.Identifier(), b.predecessor.Identifier(), "The successor's predecessor should be set after stabilizing")
}

func TestStabilizeNewSuccessor(t *testing.T) {
	a := CreateNode(1)
	b := CreateNode(2)
	c := CreateNode(4)

	a.Join(c)
	a.stabilize()

	b.Join(c)
	b.stabilize()

	assert.Equal(t, c.Identifier(), a.successor.Identifier(), "A's successor should still be C")

	// Another round of stabilization after b has joined
	a.stabilize()
	assert.Equal(t, b.Identifier(), a.successor.Identifier(), "A's successor should be B, not C")
}
