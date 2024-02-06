package chord

// TODO Make configurable
const r = 16

type SuccessorList struct {
	successors [r]node
}

// Adopt copies all values from another successor list but retains the head
func (s *SuccessorList) Adopt(t SuccessorList) {
	for i := 0; i < r-1; i++ {
		s.successors[i+1] = t.successors[i]
	}
}

// Head returns the immediate successor
func (s *SuccessorList) Head() node {
	return s.successors[0]
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
	for i := 0; i < r-2; i++ {
		if s.successors[i] == nil {
			continue
		}
		for j := i + 1; j < r-1; j++ {
			if s.successors[j] == nil {
				continue
			}
			for k := j + 1; k < r; k++ {
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
