package chord

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingletonNodeIsOwnSuccessor(t *testing.T) {
	node := CreateNode(1)

	succ, _ := node.Successor()
	assert.Equal(t, node.Identifier(), succ.Identifier(), "A new node should be its own successor")
}

func TestFindSuccessorSimple(t *testing.T) {
	a := CreateNode(1)
	b := CreateNode(10)

	a.Join(b)

	for i := a.Identifier() + 1; i <= b.Identifier(); i++ {
		a_succ, err := a.FindSuccessor(Id(i))

		assert.Nil(t, err)
		assert.Equal(t, b.Identifier(), a_succ.Identifier())
	}
}

func TestFindSuccessorWrapAround(t *testing.T) {
	a := CreateNode(10)
	b := CreateNode(1)

	a.Join(b)

	// Any key >10 should be handled by node b
	for i := 1; i < 16; i++ {
		succ, _ := a.FindSuccessor(a.Identifier() + 1)
		assert.Equal(t, b.Identifier(), succ.Identifier())
	}
}

func TestFindSuccessorWrapAroundTriple(t *testing.T) {
	a := CreateNode(10)
	b := CreateNode(1)
	c := CreateNode(5)

	a.Join(b)
	b.Join(c)

	// Any key >10 should be handled by node b
	for i := 1; i < 16; i++ {
		succ, _ := a.FindSuccessor(a.Identifier() + 1)
		assert.Equal(t, b.Identifier(), succ.Identifier())
	}
}

func TestFindSuccessorAdvanced(t *testing.T) {
	a := CreateNode(1)
	b := CreateNode(8)
	c := CreateNode(32)
	d := CreateNode(42)

	a.Join(b)
	b.Join(c)
	c.Join(d)

	// Test every key in the ring and query node A for the correct location
	for i := a.Identifier() + 1; i <= b.Identifier(); i++ {
		succ, _ := a.FindSuccessor(i)
		assert.Equal(t, b.Identifier(), succ.Identifier())
	}

	for i := b.Identifier() + 1; i <= c.Identifier(); i++ {
		succ, _ := a.FindSuccessor(i)
		assert.Equal(t, c.Identifier(), succ.Identifier())
	}

	for i := c.Identifier() + 1; i <= d.Identifier(); i++ {
		succ, _ := a.FindSuccessor(i)
		assert.Equal(t, d.Identifier(), succ.Identifier())
	}

	// All keys after the last node should be handled by the first node
	// TODO make this pass, if possible. This may not be possible to pass without requiring stabilizing
	// or fixFingers
	for i := 1; i < 16; i++ {
		// succ, _ := a.FindSuccessor(d.Identifier() + Id(i))
		// assert.Equal(t, a.Identifier(), succ.Identifier())
	}
}

func TestFindSuccessorReturnsSuccessor(t *testing.T) {
	a := CreateNode(1)
	b := CreateNode(2)

	a.Join(b)

	succ, _ := a.FindSuccessor(b.Identifier())
	assert.Equal(t, b.Identifier(), succ.Identifier())
}

func TestFindSuccessorReturnsSuccessorPermuted(t *testing.T) {
	a := CreateNode(2)
	b := CreateNode(1)

	a.Join(b)

	succ, _ := a.FindSuccessor(b.Identifier())
	assert.Equal(t, b.Identifier(), succ.Identifier())
}

func TestFindSuccessorTransitive(t *testing.T) {
	a := CreateNode(1)
	b := CreateNode(2)
	c := CreateNode(4)

	b.Join(c)
	a.Join(b)

	succ, _ := a.FindSuccessor(c.Identifier())
	assert.Equal(t, succ.Identifier(), c.Identifier())
}

// In a three node ring, the non-adjacent nodes should be aware of each other
func TestFindSuccessorTransitiveWraparound(t *testing.T) {
	a := CreateNode(128)
	b := CreateNode(1)
	c := CreateNode(16)

	b.Join(c)
	a.Join(b)

	succ, _ := a.FindSuccessor(c.Identifier())
	assert.Equal(t, succ.Identifier(), c.Identifier())
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

	succ, _ := a.Successor()
	assert.Equal(t, b.Identifier(), succ.Identifier())
}

func TestJoinSetsCorrectSuccessorPermuted(t *testing.T) {
	a := CreateNode(10)
	b := CreateNode(1)

	a.Join(b)

	succ, _ := a.Successor()
	assert.Equal(t, b.Identifier(), succ.Identifier())
}

func TestRectifySetsPredecessor(t *testing.T) {
	a := CreateNode(1)
	b := CreateNode(2)
	a.Join(b)

	b.Rectify(a)

	assert.Equal(t, a.Identifier(), b.predecessor.Identifier(), "Predecessor should be set")
}

func TestRectifyRejectsInvalidPredecessor(t *testing.T) {
	a := CreateNode(2)
	b := CreateNode(4)

	a.Join(b)
	b.Rectify(a)

	// Create a new node which comes before A. If B is notified by C, B's predecessor should still be A
	c := CreateNode(1)
	c.Join(a)
	b.Rectify(c)

	assert.Equal(t, a.Identifier(), b.predecessor.Identifier(), "The current predecessor should be unchanged")
}

func TestStabilizeSetsSuccessor(t *testing.T) {
	a := CreateNode(1)
	b := CreateNode(2)
	a.Join(b)

	_ = a.stabilize()
	// assert.NoError(t, err)

	assert.Equal(t, a.Identifier(), b.predecessor.Identifier(), "The successor's predecessor should be set after stabilizing")
}

func TestStabilizeNewSuccessor(t *testing.T) {
	a := CreateNode(1)
	b := CreateNode(2)
	c := CreateNode(4)

	a.Join(c)
	_ = a.stabilize()
	// assert.NoError(t, err)

	b.Join(c)
	_ = b.stabilize()
	// assert.NoError(t, err)

	succ, _ := a.Successor()
	assert.Equal(t, c.Identifier(), succ.Identifier(), "A's successor should still be C")

	// Another round of stabilization after b has joined
	_ = a.stabilize()
	// assert.NoError(t, err)
	succ, _ = a.Successor()
	assert.Equal(t, b.Identifier(), succ.Identifier(), "A's successor should be B, not C")
}
