package gs

import (
	"strconv"
)

func IsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func ToInt(s string) int {
	i, err := strconv.Atoi(s)
	c(err)
	return i
}
