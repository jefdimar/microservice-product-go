package common

import "strconv"

func ParseInt(str string) int {
	val, _ := strconv.Atoi(str)
	return val
}
