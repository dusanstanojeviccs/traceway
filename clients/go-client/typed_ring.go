package traceway

type TypedRing[T *E, E any] struct {
	arr      []T
	head     int
	capacity int
	len      int
}

func InitTypedRing[T *E, E any](capacity int) TypedRing[T, E] {
	return TypedRing[T, E]{
		arr:      make([]T, capacity),
		capacity: capacity,
	}
}

func (t *TypedRing[T, E]) Push(val T) {
	t.arr[t.head] = val
	t.head = (t.head + 1) % t.capacity
	if t.len < t.capacity {
		t.len += 1
	}
}

func (t *TypedRing[T, E]) ReadAll() []T {
	result := make([]T, t.len)
	for i := 0; i < t.len; i++ {
		idx := (t.head - t.len + i + t.capacity) % t.capacity
		result[i] = t.arr[idx]
	}
	return result
}

func (t *TypedRing[T, E]) Clear() {
	for i := range t.arr {
		t.arr[i] = nil
	}
	t.head = 0
	t.len = 0
}

func (t *TypedRing[T, E]) Remove(vals []T) int {
	if len(vals) == 0 {
		return 0
	}

	toRemove := make(map[T]struct{}, len(vals))
	for _, v := range vals {
		toRemove[v] = struct{}{}
	}

	writeIdx := 0
	removed := 0
	for i := 0; i < t.len; i++ {
		readIdx := (t.head - t.len + i + t.capacity) % t.capacity
		if _, shouldRemove := toRemove[t.arr[readIdx]]; shouldRemove {
			removed++
		} else {
			if writeIdx != i {
				destIdx := (t.head - t.len + writeIdx + t.capacity) % t.capacity
				t.arr[destIdx] = t.arr[readIdx]
			}
			writeIdx++
		}
	}

	for i := writeIdx; i < t.len; i++ {
		idx := (t.head - t.len + i + t.capacity) % t.capacity
		t.arr[idx] = nil
	}

	t.len = writeIdx
	t.head = (t.head - removed + t.capacity) % t.capacity

	return removed
}
