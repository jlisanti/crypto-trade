package utilities

import (
	"fmt"
	"time"
)

// specify length of moving average in hours
type MovingAverage struct {
	Length       float64
	Value        []float64
	Time         []string
	ValueSum     float64
	AverageValue float64
	Populated    bool
	NumValues    int
}

func NewMovingAverage(length float64) *MovingAverage {
	ma := MovingAverage{Length: length, NumValues: 0, ValueSum: 0.0, Populated: false}
	return &ma
}

// Starts off with zero values - wrong average until full
// Should fix this

func UpdateValue(ma *MovingAverage, newValue float64, newTime string) {
	// Determine if the moving average already contains the correct length of time
	if len(ma.Time) > 1 {
		t1, _ := time.Parse(time.RFC3339, ma.Time[0])
		t2, _ := time.Parse(time.RFC3339, newTime)
		timeDiff := t2.Sub(t1)
		fmt.Println(timeDiff)
		fmt.Println(ma.Populated)
		if (timeDiff.Seconds() / 3600.0) >= ma.Length {
			if ma.Populated == false {
				ma.Populated = true
			}
			delete := 0
			for i, _ := range ma.Value {
				t3, _ := time.Parse(time.RFC3339, ma.Time[i])
				timeDiff2 := t2.Sub(t3)

				if (timeDiff2.Seconds() / 3600.0) > ma.Length {
					delete += 1
					ma.ValueSum -= ma.Value[i]
				} else {
					break
				}
			}
			ma.ValueSum += newValue
			ma.Value = append(ma.Value, newValue)
			ma.Value = ma.Value[delete:]
			ma.Time = append(ma.Time, newTime)
			ma.Time = ma.Time[delete:]
			ma.AverageValue = ma.ValueSum / float64(len(ma.Value)-1)
			fmt.Println(ma.ValueSum)
			fmt.Println(len(ma.Value))
		} else {
			ma.Value = append(ma.Value, newValue)
			ma.Time = append(ma.Time, newTime)
			ma.ValueSum += newValue
			ma.AverageValue = ma.ValueSum / float64(len(ma.Value)-1)
		}
	} else {
		ma.Value = append(ma.Value, newValue)
		ma.Time = append(ma.Time, newTime)
		ma.ValueSum += newValue
		ma.AverageValue = ma.ValueSum // float64(len(ma.Value))

	}
}
