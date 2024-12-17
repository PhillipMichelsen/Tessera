package component

import (
	"AlgorithimcTraderDistributed/common/models"
)

type MovingAverageOHLCVRollingComponent struct {
	Period     int
	data       []models.OHLCV
	currentSum float64
}

func NewMovingAverageOHLCVRollingComponent(period int) *MovingAverageOHLCVRollingComponent {
	return &MovingAverageOHLCVRollingComponent{
		Period:     period,
		data:       []models.OHLCV{},
		currentSum: 0,
	}
}

func (m *MovingAverageOHLCVRollingComponent) AddNewData(newData models.OHLCV) {
	m.data = append(m.data, newData)
	m.currentSum += newData.Close
	if len(m.data) > m.Period {
		m.currentSum -= m.data[0].Close
		m.data = m.data[1:]
	}
}

func (m *MovingAverageOHLCVRollingComponent) GetMovingAverage() float64 {
	return m.currentSum / float64(len(m.data))
}

func (m *MovingAverageOHLCVRollingComponent) IsReady() bool {
	return len(m.data) == m.Period
}
