package testrepo

import "fmt"

// Calculator performs basic arithmetic
type Calculator struct {
	Precision int
}

// Add sums two integers and returns the result
func (c *Calculator) Add(a, b int) int {
	return a + b
}

// FormatResult takes a value and turns it into a pretty string
func FormatResult(val int) string {
	return fmt.Sprintf("The result is: %d", val)
}
