package tack

type entry struct {
	hit int64
	value string
}

// Calculates the total memory usage including the hit timestamp.
func (d entry) getMem() int {
	return len(d.value) + 8
}

func (d entry) getValueMem() int {
	return len(d.value)
}