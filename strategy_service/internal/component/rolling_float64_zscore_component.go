package component

import "math"

type RollingFloat64ZScoreComponent struct {
	Period int
	data   []float64
}

func NewRollingFloat64ZScoreComponent(period int) *RollingFloat64ZScoreComponent {
	return &RollingFloat64ZScoreComponent{
		Period: period,
		data:   []float64{},
	}
}

func (r *RollingFloat64ZScoreComponent) AddNewData(newData float64) {
	r.data = append(r.data, newData)
	if len(r.data) > r.Period {
		r.data = r.data[1:]
	}
}

func (r *RollingFloat64ZScoreComponent) GetZScore() float64 {
	if len(r.data) == 0 {
		return 0.0 // Avoid calculation if no data is present
	}

	mean := r.mean()
	stdDev := r.stdDev(mean)

	if stdDev == 0 {
		return 0.0 // Avoid division by zero
	}

	return (r.data[len(r.data)-1] - mean) / stdDev
}

func (r *RollingFloat64ZScoreComponent) mean() float64 {
	sum := 0.0
	for _, val := range r.data {
		sum += val
	}
	return sum / float64(len(r.data))
}

func (r *RollingFloat64ZScoreComponent) stdDev(mean float64) float64 {
	sum := 0.0
	for _, val := range r.data {
		sum += (val - mean) * (val - mean)
	}
	variance := sum / float64(len(r.data))
	return math.Sqrt(variance) // Take square root for standard deviation
}

func (r *RollingFloat64ZScoreComponent) IsReady() bool {
	return len(r.data) == r.Period
}
