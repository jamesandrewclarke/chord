package chord

import (
	"fmt"
	"strings"
)

// TODO Make configurable
const SUCCESSOR_LIST_SIZE = 3

type SuccessorList struct {
	successors [SUCCESSOR_LIST_SIZE]node
}

// Adopt copies all values from another successor list but retains the head
func (s *SuccessorList) Adopt(t SuccessorList) {
	for i := 0; i < SUCCESSOR_LIST_SIZE-1; i++ {
		s.successors[i+1] = t.successors[i]
	}
}

// Head returns the immediate successor
func (s *SuccessorList) Head() node {
	return s.successors[0]
}

// Removes the first element of the list
func (s *SuccessorList) PopHead() {
	// Shifts all elements back one place
	for i := 1; i < SUCCESSOR_LIST_SIZE; i++ {
		s.successors[i-1] = s.successors[i]
	}
	s.successors[SUCCESSOR_LIST_SIZE-1] = nil
}

// SetHead sets the immediate successor
func (s *SuccessorList) SetHead(p node) {
	s.successors[0] = p
}

// Ordered checks the 'EvaluatedSuccessorList' invariant
// Intended for use outside a main loop for local monitoring
func (s *SuccessorList) Ordered() bool {
	// TODO Could be better than O(r^3)?
	// Exhaustive check for now just to be sure
	for i := 0; i < SUCCESSOR_LIST_SIZE-2; i++ {
		if s.successors[i] == nil {
			continue
		}
		for j := i + 1; j < SUCCESSOR_LIST_SIZE-1; j++ {
			if s.successors[j] == nil {
				continue
			}
			for k := j + 1; k < SUCCESSOR_LIST_SIZE; k++ {
				if s.successors[k] == nil {
					continue
				}
				if !between(s.successors[j].Identifier(), s.successors[i].Identifier(), s.successors[k].Identifier()) {
					return false
				}
			}
		}
	}

	return true
}

// UniqueSuccessors checks that the successor list contains no duplicate values
// Intended for use outside a main loop for local monitoring
func (s *SuccessorList) UniqueSuccessors() bool {
	identifiers := make(map[Id]bool)
	for _, succ := range s.successors {
		if succ == nil {
			continue
		}
		if _, ok := identifiers[succ.Identifier()]; ok {
			return false
		}
		identifiers[succ.Identifier()] = true
	}

	return true
}

func (s *SuccessorList) String() string {
	line := ""
	for _, succ := range s.successors {
		if succ != nil {
			line = strings.Join([]string{line, fmt.Sprintf("%v, ", succ.Identifier())}, "")
		} else {
			line = strings.Join([]string{line, "* "}, "")
		}
	}

	return line
}
