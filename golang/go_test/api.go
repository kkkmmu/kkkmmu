package api

import (
	"fmt"
)

func api_return(i int) int {
	return i
}

func api_increase(i int) int {
	return i + 1
}

func api_return_str(i int) string {
	return fmt.Sprintf("%d", i)
}
