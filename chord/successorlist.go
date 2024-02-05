package chord

const r = 128

type SuccessorList struct {
	successors [r]node
}

func (s *SuccessorList) Adopt(t SuccessorList) {
	for i := 0; i < r-1; i++ {
		s.successors[i+1] = t.successors[i]
	}
}

func (s *SuccessorList) Head() node {
	return s.successors[0]
}

func (s *SuccessorList) SetHead(p node) {
	s.successors[0] = p
}
