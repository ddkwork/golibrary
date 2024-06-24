package xlsx

import (
	"math"
	"time"
)

func timeFromExcelTime(excelTime float64) time.Time {
	var date time.Time
	intPart := int64(excelTime)

	if intPart <= 61 {
		return julianDateToGregorianTime(2400000.5, excelTime+15018.0)
	}
	floatPart := excelTime - float64(intPart)
	var dayNanoSeconds float64 = 24 * 60 * 60 * 1000 * 1000 * 1000
	date = time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	durationDays := time.Duration(intPart) * time.Hour * 24
	durationPart := time.Duration(dayNanoSeconds * floatPart)
	return date.Add(durationDays).Add(durationPart)
}

func julianDateToGregorianTime(part1, part2 float64) time.Time {
	part1I, part1F := math.Modf(part1)
	part2I, part2F := math.Modf(part2)
	julianDays := part1I + part2I
	julianFraction := part1F + part2F
	julianDays, julianFraction = shiftJulianToNoon(julianDays, julianFraction)
	day, month, year := fliegelAndVanFlandernAlgorithm(int(julianDays))
	hours, minutes, seconds, nanoseconds := fractionOfADay(julianFraction)
	return time.Date(year, time.Month(month), day, hours, minutes, seconds, nanoseconds, time.UTC)
}

func shiftJulianToNoon(julianDays, julianFraction float64) (julianDaysResult, julianFractionResult float64) {
	switch {
	case -0.5 < julianFraction && julianFraction < 0.5:
		julianFraction += 0.5
	case julianFraction >= 0.5:
		julianDays++
		julianFraction -= 0.5
	case julianFraction <= -0.5:
		julianDays--
		julianFraction += 1.5
	}
	return julianDays, julianFraction
}

func fliegelAndVanFlandernAlgorithm(jd int) (day, month, year int) {
	l := jd + 68569
	n := (4 * l) / 146097
	l -= (146097*n + 3) / 4
	i := (4000 * (l + 1)) / 1461001
	l = l - (1461*i)/4 + 31
	j := (80 * l) / 2447
	d := l - (2447*j)/80
	l = j / 11
	m := j + 2 - (12 * l)
	y := 100*(n-49) + i + l
	return d, m, y
}

func fractionOfADay(fraction float64) (hours, minutes, seconds, nanoseconds int) {
	const (
		c1us = 1e3
		c1s  = 1e9
	)
	frac := int64(24*60*60*c1s*fraction + c1us/2)
	nanoseconds = int((frac%c1s)/c1us) * c1us
	frac /= c1s
	seconds = int(frac % 60)
	frac /= 60
	minutes = int(frac % 60)
	hours = int(frac / 60)
	return
}
