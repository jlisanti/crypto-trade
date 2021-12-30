package utilities

import (
	"time"
)

// specify length of moving average in hours
type MovingAverage struct {
	Length       float64
	Value        []float64
	Time         []time.Time
	TimeDiffs    []float64
	TimeValues   []string
	ValueSum     float64
	AverageValue float64
	Averages     []float64
	Populated    bool
	NumValues    int
}

func NewMovingAverage(length float64) *MovingAverage {
	ma := MovingAverage{Length: length, NumValues: 0, ValueSum: 0.0, Populated: false}
	return &ma
}

func UpdateValue(ma *MovingAverage, newValue float64, newTime time.Time) {

	// Determine if the moving average already contains the correct length of time
	if len(ma.Time) > 1 {

		t1 := ma.Time[0]
		t2 := newTime
		timeDiff := t2.Sub(t1)

		if !ma.Populated {
			ma.Value = append(ma.Value, newValue)
			ma.Time = append(ma.Time, newTime)
			ma.ValueSum += newValue
			ma.TimeDiffs = append(ma.TimeDiffs, timeDiff.Hours())

			ma.AverageValue = ma.ValueSum / float64(len(ma.Value))
			ma.Averages = append(ma.Averages, ma.ValueSum/float64(len(ma.Value)))
			ma.TimeValues = append(ma.TimeValues, newTime.String())

			if timeDiff.Hours() >= ma.Length {
				ma.Populated = true
			}

		} else {

			delete := 0
			for i, _ := range ma.Value {
				timeDiff2 := newTime.Sub(ma.Time[i])
				if timeDiff2.Hours() > ma.Length {
					delete += 1
					ma.ValueSum -= ma.Value[i]
				} else {
					break
				}
			}

			ma.ValueSum += newValue

			ma.Value = append(ma.Value, newValue)
			ma.Time = append(ma.Time, newTime)
			ma.Averages = append(ma.Averages, newValue)
			ma.TimeValues = append(ma.TimeValues, newTime.String())
			ma.TimeDiffs = append(ma.TimeDiffs, timeDiff.Hours())

			if delete != 0 {
				ma.Value = append(ma.Value[:delete], ma.Value[delete+1:]...)
				ma.Time = append(ma.Time[:delete], ma.Time[delete+1:]...)
				ma.Averages = append(ma.Averages[:delete], ma.Averages[delete+1:]...)
				ma.TimeValues = append(ma.TimeValues[:delete], ma.TimeValues[delete+1:]...)
				ma.TimeDiffs = append(ma.TimeDiffs[:delete], ma.TimeDiffs[delete+1:]...)
			}

			ma.AverageValue = ma.ValueSum / float64(len(ma.Value))
		}

	} else {
		ma.Value = append(ma.Value, newValue)
		ma.Time = append(ma.Time, newTime)
		ma.ValueSum += newValue
		ma.AverageValue = ma.ValueSum / float64(len(ma.Value))
		ma.Averages = append(ma.Averages, ma.ValueSum/float64(len(ma.Value)))
		ma.TimeValues = append(ma.TimeValues, newTime.String())
		ma.TimeDiffs = append(ma.TimeDiffs, 0.0)
	}
}
