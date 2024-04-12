package chord

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBetweenReturnsTrueForCorrectRange(t *testing.T) {
	ranges := [][3]Id{
		{0, 100, 20},
		{5, 0, 10},
		{1, 0, 2},
	}

	for _, v := range ranges {
		assert.True(t, Between(v[0], v[1], v[2]))
	}
}

func TestBetweenReturnsFalseForIncorrectRange(t *testing.T) {
	ranges := [][3]Id{
		{1, 1, 1},
		{10, 10, 10},
		{10, 0, 5},
		{90, 100, 20},
	}

	for _, v := range ranges {
		assert.False(t, Between(v[0], v[1], v[2]))
	}
}
