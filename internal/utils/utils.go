package utils

import (
	"fmt"
	"strconv"
)

func ConvertStringToUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to uint: %w", err)
	}
	return uint(val), nil
}
