package chord

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingletonNodeIsOwnSuccessor(t *testing.T) {
	node := CreateNode(1)

	assert.Equal(t, node.Identifier(), node.successor.Identifier(), "A new node should be its own successor")
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
