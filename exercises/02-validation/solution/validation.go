// SOLUTION FILE - For mentors only!
package validation

import (
	"errors"
	"fmt"
	"math"
)

type GLEntry struct {
	Account string
	Debit   float64
	Credit  float64
}

const epsilon = 0.0001

func TotalDebit(entries []GLEntry) float64 {
	var total float64
	for _, entry := range entries {
		total += entry.Debit
	}
	return total
}

func TotalCredit(entries []GLEntry) float64 {
	var total float64
	for _, entry := range entries {
		total += entry.Credit
	}
	return total
}

func IsBalanced(entries []GLEntry) bool {
	debit := TotalDebit(entries)
	credit := TotalCredit(entries)
	return math.Abs(debit-credit) < epsilon
}

func Difference(entries []GLEntry) float64 {
	return TotalDebit(entries) - TotalCredit(entries)
}

func ValidateGLMap(entries []GLEntry) error {
	if len(entries) == 0 {
		return errors.New("GL map cannot be empty")
	}

	for i, entry := range entries {
		if entry.Account == "" {
			return fmt.Errorf("entry %d: account is required", i)
		}
	}

	if !IsBalanced(entries) {
		diff := Difference(entries)
		return fmt.Errorf("entries not balanced: difference = %.2f", diff)
	}

	return nil
}
