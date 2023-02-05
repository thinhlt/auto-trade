package utils

type QueueInt64 struct {
	elements []int64
	max      int
}

func NewQueue() QueueInt64 {
	return QueueInt64{
		elements: make([]int64, 0, 3),
		max:      3,
	}
}

func NewQueueLimit(max int) QueueInt64 {
	return QueueInt64{
		elements: make([]int64, 0, max),
		max:      max,
	}
}

func NewQueueWithData(max int, data []int64) QueueInt64 {
	return QueueInt64{
		elements: data,
		max:      max,
	}
}

func (q *QueueInt64) EnQueue(f int64) {
	if q.max == 0 {
		q.elements = []int64{f}
		q.max = 1
		return
	}
	if len(q.elements) == q.max {
		q.elements[0] = 0
		q.elements = append(q.elements[1:], f)
	} else {
		q.elements = append(q.elements, f)
	}
}

func (q *QueueInt64) DeQueue() int64 {
	if len(q.elements) == 0 {
		return 0
	}
	q.elements = q.elements[1:]
	result := q.elements[0]
	q.elements[0] = 0
	return result
}

func (q QueueInt64) GetSlice() []int64 {
	return q.elements
}

func (q *QueueInt64) GetDeltaAvg() float64 {
	length := len(q.elements)
	if length == 0 {
		return 0
	}
	sum := int64(0)
	for _, v := range q.elements {
		sum += v
	}
	return float64(sum) * 1.0 / float64(length)
}

func (q QueueInt64) GetLength() int {
	return len(q.elements)
}

func (q QueueInt64) GetLast() int64 {
	return q.elements[q.GetLength()-1]
}

func (q QueueInt64) GetLast3Avg() float64 {
	length := q.GetLength()
	if length == 0 {
		return 0
	}
	var tmp []int64
	if length < 3 {
		tmp = q.elements
	} else {
		tmp = q.elements[length-3:]
	}
	sum := int64(0)
	for _, v := range tmp {
		sum += v
	}
	return float64(sum) * 1.0 / float64(length)
}
