package utils

import (
	"fmt"
	"strconv"
	"time"
)

func ConvertStringToUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to uint: %w", err)
	}
	return uint(val), nil
}

func GetWeekDateRange(date time.Time) WeeklyDateRange {
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	startOfWeek := date.AddDate(0, 0, -weekday+1)
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, loc)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)
	endOfWeek = time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 0, loc)

	return WeeklyDateRange{
		StartOfWeek: startOfWeek,
		EndOfWeek:   endOfWeek,
	}
}

func GenerateWeeklyDateRanges(start time.Time, n int) []WeeklyDateRange {
	var ranges []WeeklyDateRange
	curr := GetWeekDateRange(start)
	loc, _ := time.LoadLocation("Asia/Jakarta")
	for range n + 1 {
		ranges = append(ranges, WeeklyDateRange{
			StartOfWeek: curr.StartOfWeek,
			EndOfWeek:   curr.EndOfWeek,
		})
		curr.StartOfWeek = curr.StartOfWeek.AddDate(0, 0, 7)
		curr.StartOfWeek = time.Date(curr.StartOfWeek.Year(), curr.StartOfWeek.Month(), curr.StartOfWeek.Day(), 0, 0, 0, 0, loc)
		curr.EndOfWeek = curr.EndOfWeek.AddDate(0, 0, 7)
		curr.EndOfWeek = time.Date(curr.EndOfWeek.Year(), curr.EndOfWeek.Month(), curr.EndOfWeek.Day(), 23, 59, 59, 0, loc)
	}
	return ranges
}
