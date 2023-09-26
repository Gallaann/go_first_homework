package stack

type Stack []string

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(str string) {
	*s = append(*s, str)
}

func (s *Stack) Pop() string {
	if s.IsEmpty() {
		return ""
	}

	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]

	return element
}

func (s *Stack) Peek() string {
	if s.IsEmpty() {
		return ""
	}

	return (*s)[len(*s)-1]
}
