package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Step struct {
	MinUnit int
	Rate    float32
}

type UsageStep struct {
	Units int
	Step  *Step
	Cost  float32
}

type UsageExtra struct {
	Name string
	Cost float32
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func calculateBase(steps []Step, units int) ([]UsageStep, float32) {
	var l int = len(steps) - 1
	var cost float32 = 0
	var calculatedUnits int = 0
	var usages []UsageStep
	for i := 0; i <= l; i++ {
		remainingUnit := units - calculatedUnits
		stepUnit := 0

		if remainingUnit == 0 {
			break
		}

		if i < l {
			stepUnit = Min(steps[i+1].MinUnit-steps[i].MinUnit, remainingUnit)
		} else {
			stepUnit = remainingUnit
		}

		usage := UsageStep{}
		usage.Step = &steps[i]
		usage.Units = stepUnit
		usage.Cost = usage.Step.Rate * float32(usage.Units)
		usages = append(usages, usage)

		calculatedUnits += stepUnit
		cost += usage.Cost
	}

	return usages, cost
}

func calculate(units int) ([]UsageStep, []UsageExtra, float32) {
	stepsPdbSep2015 := []Step{
		{1, 3.80},
		{76, 5.14},
		{201, 5.36},
		{301, 5.63},
		{401, 8.70},
		{601, 9.98}}
	usages, baseCost := calculateBase(stepsPdbSep2015, units)

	var extras []UsageExtra

	extras = append(extras, UsageExtra{"Demand Charge", 45})
	extras = append(extras, UsageExtra{"Service Charge", 10})
	extras = append(extras, UsageExtra{"VAT (5.0025%)", (baseCost + 45 + 10) * 5.0025 / 100})
	return usages, extras, baseCost
}

func main() {

	router := gin.Default()
	router.GET("/calculate/:unit", func(c *gin.Context) {
		unit, err := strconv.Atoi(c.Param("unit"))
		if err != nil {
			c.AbortWithStatus(422)
		}
		usages, extras, cost := calculate(unit)
		c.JSON(200, gin.H{
			"usages":   usages,
			"extras":   extras,
			"baseCost": cost,
		})
	})

	router.Run(":8011")
}
