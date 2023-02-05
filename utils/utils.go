package utils

func AverageFloat(values ...float64) float64 {
	s := float64(0)
	n := float64(len(values))
	for _, v := range values {
		s += v
	}
	result := s / n
	return result
}

func AverageInt(values ...int) float64 {
	s := 0
	n := float64(len(values))
	for _, v := range values {
		s += v
	}
	result := float64(s) / n
	return result
}

func FindMaxFloat(values ...float64) float64 {
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

func FindMinFloat(values ...float64) float64 {
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func FindMaxNMinFloat(values ...float64) (max, min float64) {
	max = values[0]
	min = values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return max, min
}
