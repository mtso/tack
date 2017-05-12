package tack

type entry struct {
	hit int64
	value string
}

// Calculates the total memory usage including the hit timestamp.
func (e *entry) getMem() int {
	return len(e.value) + 8
}

func (e *entry) getValueMem() int {
	return len(e.value)
}

func (e *entry) setHit(nsec int64) {
	e.hit = nsec
}