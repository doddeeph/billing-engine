package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvertStringToUint(t *testing.T) {
	actual, _ := ConvertStringToUint("1")
	assert.Equal(t, uint(1), actual)

	actual, err := ConvertStringToUint("one")
	if assert.Error(t, err) {
		assert.Equal(t, "failed to convert string to uint: strconv.ParseUint: parsing \"one\": invalid syntax", err.Error())
	}
}

func TestGetWeekRange(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	today := time.Date(2025, 8, 7, 8, 30, 30, 30, loc)
	actual := GetWeekDateRange(today)

	expectedStartOfWeek := time.Date(2025, 8, 4, 0, 0, 0, 0, loc)
	assert.Equal(t, expectedStartOfWeek, actual.StartOfWeek)
	expectedEndOfWeek := time.Date(2025, 8, 10, 23, 59, 59, 0, loc)
	assert.Equal(t, expectedEndOfWeek, actual.EndOfWeek)
}

func TestGenerateWeeklyDateRanges(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	today := time.Date(2025, 8, 7, 8, 30, 30, 30, loc)
	actual := GenerateWeeklyDateRanges(today, 4)

	expectedStartOfWeek := time.Date(2025, 8, 4, 0, 0, 0, 0, loc)
	assert.Equal(t, expectedStartOfWeek, actual[0].StartOfWeek)
	expectedEndOfWeek := time.Date(2025, 8, 10, 23, 59, 59, 0, loc)
	assert.Equal(t, expectedEndOfWeek, actual[0].EndOfWeek)

	expectedStartOfWeek = time.Date(2025, 8, 11, 0, 0, 0, 0, loc)
	assert.Equal(t, expectedStartOfWeek, actual[1].StartOfWeek)
	expectedEndOfWeek = time.Date(2025, 8, 17, 23, 59, 59, 0, loc)
	assert.Equal(t, expectedEndOfWeek, actual[1].EndOfWeek)

	expectedStartOfWeek = time.Date(2025, 9, 1, 0, 0, 0, 0, loc)
	assert.Equal(t, expectedStartOfWeek, actual[4].StartOfWeek)
	expectedEndOfWeek = time.Date(2025, 9, 7, 23, 59, 59, 0, loc)
	assert.Equal(t, expectedEndOfWeek, actual[4].EndOfWeek)
}
