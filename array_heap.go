package multitimer

type arrayHeap[T any, Less func(a, b T) bool] struct {
	less Less
	data []T
}

func newArrayHeap[T any, Less func(a, b T) bool](less Less) *arrayHeap[T, Less] {
	return &arrayHeap[T, Less]{
		less: less,
	}
}

func (ah *arrayHeap[T, Less]) Len() int {
	return len(ah.data)
}

func (ah *arrayHeap[T, Less]) Less(i, j int) bool {
	return ah.less(ah.data[i], ah.data[j])
}

func (ah *arrayHeap[T, Less]) Swap(i, j int) {
	ah.data[i], ah.data[j] = ah.data[j], ah.data[i]
}

func (ah *arrayHeap[T, Less]) Push(x any) {
	ah.data = append(ah.data, x.(T))
}

func (ah *arrayHeap[T, less]) Top() T {
	return ah.data[0]
}

func (ah *arrayHeap[T, Less]) Pop() any {
	last := len(ah.data) - 1
	ret := ah.data[last]
	ah.data = ah.data[:last]
	return ret
}

func (ah *arrayHeap[T, less]) Clear() {
	ah.data = nil
}
