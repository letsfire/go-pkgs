package utils

// UniqueStringSlice
type UniqueStringSlice struct {
	data []string
	hash map[string]int
}

func (uss *UniqueStringSlice) Append(vs ...string) {
	if uss.hash == nil {
		uss.hash = make(map[string]int)
	}
	for _, v := range vs {
		uss.hash[v] += 1
		if uss.hash[v] == 1 {
			uss.data = append(uss.data, v)
		}
	}
}

func (uss *UniqueStringSlice) GetValue() []string {
	if uss.data == nil {
		uss.data = make([]string, 0)
	}
	return uss.data
}

func (uss *UniqueStringSlice) GetStats() map[string]int {
	if uss.hash == nil {
		uss.hash = make(map[string]int)
	}
	return uss.hash
}
