package chord

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadReturnsCorrectNode(t *testing.T) {
	s := SuccessorList{}
	n := CreateNode(1)

	s.SetHead(n)

	assert.Equal(t, n, s.Head())
}

func TestAdoptSimple(t *testing.T) {
	s := SuccessorList{}
	testNode := CreateNode(1)
	s.SetHead(testNode)

	u := SuccessorList{}
	testNode2 := CreateNode(2)
	u.SetHead(testNode2)

	s.Adopt(u)

	assert.Equal(t, testNode, s.successors[0])
	assert.Equal(t, testNode2, s.successors[1])
}

func TestAdoptAdvanced(t *testing.T) {
	s := SuccessorList{}
	for i := 0; i < r; i++ {
		s.successors[i] = CreateNode(Id(i))
	}

	u := SuccessorList{}
	for i := 0; i < r; i++ {
		u.successors[i] = CreateNode(Id(i + 1000))
	}

	s.Adopt(u)

	assert.Equal(t, Id(0), s.Head().Identifier())
	for i := 1; i < r; i++ {
		assert.Equal(t, Id(i+1000-1), s.successors[i].Identifier())
	}
}
