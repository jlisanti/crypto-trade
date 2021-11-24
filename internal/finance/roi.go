package finance

import (
	"math"
	"strconv"
	"time"
)

func ComputeROI(sellValue string, quantity string, buyPrice string, fees string, buyDate time.Time) (roi float64, value float64, age time.Duration) {

	finalPriceOfInvestment, _ := strconv.ParseFloat(sellValue, 64)
	quantityOfInvestment, _ := strconv.ParseFloat(quantity, 64)
	finalValueOfInvestment := finalPriceOfInvestment * quantityOfInvestment

	initialValueOfInvestment, _ := strconv.ParseFloat(buyPrice, 64)
	initialValueOfInvestment = math.Abs(initialValueOfInvestment)

	costOfInvestment, _ := strconv.ParseFloat(fees, 64)

	roi = ((finalValueOfInvestment - initialValueOfInvestment) / ((2.0 * costOfInvestment) + initialValueOfInvestment)) * 100.0
	age = time.Now().Sub(buyDate)
	value = finalValueOfInvestment - initialValueOfInvestment

	return roi, value, age
}
