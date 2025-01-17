package chord

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadReturnsCorrectNode(t *testing.T) {
	s := CreateSuccessorList(10)
	n := CreateNode(1)

	s.SetHead(n)

	assert.Equal(t, n, s.Head())
}

func TestPopHeadRemovesCorrectNode(t *testing.T) {
	s := CreateSuccessorList(10)
	for i := 0; i < 10; i++ {
		s.successors[i] = CreateNode(Id(i))
	}
	s.PopHead()
	assert.Equal(t, s.Head().Identifier(), Id(1))
	for i := 0; i < 9; i++ {
		assert.Equal(t, s.successors[i].Identifier(), Id(i+1))
	}
	assert.Nil(t, s.successors[9], "final element should be nil")
}

func TestAdoptSimple(t *testing.T) {
	s := CreateSuccessorList(10)
	testNode := CreateNode(1)
	s.SetHead(testNode)

	u := CreateSuccessorList(10)
	testNode2 := CreateNode(2)
	u.SetHead(testNode2)

	s.Adopt(u)

	assert.Equal(t, testNode, s.successors[0])
	assert.Equal(t, testNode2, s.successors[1])
}

func TestAdoptAdvanced(t *testing.T) {
	s := CreateSuccessorList(10)
	for i := 0; i < 10; i++ {
		s.successors[i] = CreateNode(Id(i))
	}

	u := CreateSuccessorList(10)
	for i := 0; i < 10; i++ {
		u.successors[i] = CreateNode(Id(i + 1000))
	}

	s.Adopt(u)

	assert.Equal(t, Id(0), s.Head().Identifier())
	for i := 1; i < 10; i++ {
		assert.Equal(t, Id(i+1000-1), s.successors[i].Identifier())
	}
}

func TestUniqueSuccessorsFalse(t *testing.T) {
	s := CreateSuccessorList(10)
	s.successors[0] = CreateNode(0)
	s.successors[1] = CreateNode(0)

	assert.False(t, s.UniqueSuccessors())
}

func TestUniqueSuccessorsTrue(t *testing.T) {
	s := CreateSuccessorList(10)

	for i := 0; i < 10; i++ {
		s.successors[i] = CreateNode(Id(i))
	}

	assert.True(t, s.UniqueSuccessors())
}

func TestUniqueSuccessorsEmpty(t *testing.T) {
	s := CreateSuccessorList(10)
	assert.True(t, s.UniqueSuccessors())
}

func TestOrderedTrue(t *testing.T) {
	table := [][]Id{
		// Ascending
		{1, 2, 3},
		{1, 2, 3, 4, 5, 6, 7},
		{1, 16, 256, 4096, 65536},

		// Wrap-around
		{2, 3, 1},
		{4, 5, 1, 2, 3},
		{100, 200, 500, 50},

		// Edge cases
		{1},
		{1, 2},
		{},
	}

	for _, nums := range table {
		s := CreateSuccessorList(10)
		for i, num := range nums {
			s.successors[i] = CreateNode(num)
		}
		assert.True(t, s.Ordered(), "%v should be true", nums)
	}
}

func TestOrderedFalse(t *testing.T) {
	table := [][]Id{
		{4, 2, 6},
		{1, 10, 2},
		{100, 104, 102, 105, 107, 250},
		{100, 104, 105, 106, 103},
		{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
	}

	for _, nums := range table {
		s := CreateSuccessorList(10)
		for i, num := range nums {
			s.successors[i] = CreateNode(num)
		}
		assert.False(t, s.Ordered(), "%v should be false", nums)
	}
}
