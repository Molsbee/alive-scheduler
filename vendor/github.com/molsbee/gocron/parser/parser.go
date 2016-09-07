package parser

import (
	"log"
	"math"
	"strconv"
	"strings"
)

var (
	monthNames = map[string]uint{
		"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6, "jul": 7,
		"aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
	}
	dayNames = map[string]uint{
		"sun": 1, "mon": 2, "tue": 3, "wed": 4, "thu": 5, "fri": 6, "sat": 7,
	}
)

// Parse converts a classic cron formatted string into a schedule which
// provides mechanism to determine the next time the task should be executed.
// * * * * *
// | | | | |
// | | | | |
// | | | | +---- Day of the Week   (range: 1-7, 1 standing for Monday)
// | | | +------ Month of the Year (range: 1-12)
// | | +-------- Day of the Month  (range: 1-31)
// | +---------- Hour              (range: 0-23)
// +------------ Minute            (range: 0-59)
func Parse(cron string) *Schedule {
	fields := strings.Fields(cron)
	if len(fields) != 5 {
		log.Panic("Incorrect schedule provided")
	}

	return &Schedule{
		Minute:     getField(fields[0], 0, 59, nil),
		Hour:       getField(fields[1], 0, 23, nil),
		DayOfMonth: getField(fields[2], 1, 31, nil),
		Month:      getField(fields[3], 1, 12, monthNames),
		DayOfWeek:  getField(fields[4], 1, 7, dayNames),
	}
}

func getField(v string, min, max uint, names map[string]uint) uint64 {
	values := strings.FieldsFunc(v, func(r rune) bool {
		return r == ','
	})

	if len(values) != 1 {
		var bits uint64
		for _, value := range values {
			bits |= getField(value, min, max, names)
		}
		return bits
	}

	stepValues := strings.Split(v, "/")
	step := calculateStep(stepValues) // Validate & Calculate number of slashes

	rangeValues := strings.Split(stepValues[0], "-")
	if len(rangeValues) < 0 || len(rangeValues) > 2 {
		log.Panicf("Too many hyphens provided: %s", v)
	}

	if rangeValues[0] == "*" || rangeValues[0] == "?" {
		return getBits(min, max, step)
	}

	start := parseField(rangeValues[0], names)
	end := start
	if start < min {
		log.Panicf("Value (%d) provided is below the minimum allowed value (%d) for field", start, min)
	}

	if len(rangeValues) == 2 {
		end = parseField(rangeValues[1], names)
		if end > max {
			log.Panicf("Value (%d) provided is above the maximum allowed value (%d) for field", end, max)
		}
		if start > end {
			log.Panicf("Beginning of range (%d) is beyond the end range (%d)", start, end)
		}
	}

	return getBits(start, end, step)
}

func calculateStep(stepValues []string) (step uint) {
	switch len(stepValues) {
	case 1:
		step = 1
	case 2:
		step = parseField(stepValues[1], nil)
	default:
		log.Panicf("To many (/) provided - %s", stepValues)
	}
	return
}

func parseField(field string, names map[string]uint) uint {
	if names != nil {
		if value, ok := names[strings.ToLower(field)]; ok {
			return value
		}
	}
	value, err := strconv.Atoi(field)
	if err != nil || value < 0 {
		log.Panicf("Failed to parse field or incorrect value provided %s", field)
	}

	return uint(value)
}

func getBits(min, max, step uint) uint64 {
	if step == 1 {
		return ^(math.MaxUint64 << (max + 1)) & (math.MaxUint64 << min)
	}

	var bits uint64
	for i := min; i <= max; i += step {
		bits |= 1 << i
	}
	return bits
}
