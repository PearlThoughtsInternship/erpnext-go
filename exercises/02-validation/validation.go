package validation

import (
	"math"
)

// GLEntry represents a General Ledger Entry.
// (Simplified version for this exercise)
type GLEntry struct {
	Account string
	Debit   float64
	Credit  float64
}

// epsilon is the tolerance for floating point comparisons.
// Due to how computers store decimals, 0.1 + 0.2 might equal 0.30000000000000004
// We use epsilon to handle these tiny differences.
const epsilon = 0.0001

// TotalDebit calculates the sum of all debit amounts in the entries.
//
// Example:
//
//	entries := []GLEntry{
//	    {Debit: 100},
//	    {Debit: 50},
//	}
//	TotalDebit(entries) // returns 150
//
// TODO: Implement this function
func TotalDebit(entries []GLEntry) float64 {
	// YOUR CODE HERE
	// Hint: Use a for loop to iterate over entries
	// and sum up the Debit field

	return 0 // Replace this
}

// TotalCredit calculates the sum of all credit amounts in the entries.
//
// Example:
//
//	entries := []GLEntry{
//	    {Credit: 100},
//	    {Credit: 50},
//	}
//	TotalCredit(entries) // returns 150
//
// TODO: Implement this function
func TotalCredit(entries []GLEntry) float64 {
	// YOUR CODE HERE
	// Hint: Same pattern as TotalDebit, but for Credit field

	return 0 // Replace this
}

// IsBalanced checks if the total debits equal total credits.
// Uses epsilon for floating point comparison.
//
// Example:
//
//	balanced := []GLEntry{
//	    {Debit: 100},
//	    {Credit: 100},
//	}
//	IsBalanced(balanced) // returns true
//
//	unbalanced := []GLEntry{
//	    {Debit: 100},
//	    {Credit: 90},
//	}
//	IsBalanced(unbalanced) // returns false
//
// TODO: Implement this function
func IsBalanced(entries []GLEntry) bool {
	// YOUR CODE HERE
	// Hint:
	// 1. Get total debit using TotalDebit()
	// 2. Get total credit using TotalCredit()
	// 3. Compare using math.Abs(debit - credit) < epsilon

	_ = math.Abs // This line is here so the import doesn't error

	return false // Replace this
}

// Difference calculates debit minus credit.
// Positive = more debits, Negative = more credits.
//
// Example:
//
//	entries := []GLEntry{
//	    {Debit: 100},
//	    {Credit: 80},
//	}
//	Difference(entries) // returns 20 (more debits)
//
// TODO: Implement this function
func Difference(entries []GLEntry) float64 {
	// YOUR CODE HERE
	// Hint: TotalDebit(entries) - TotalCredit(entries)

	return 0 // Replace this
}

// ValidateGLMap validates a slice of GL entries and returns an error if invalid.
// This is a more complete validation function that checks multiple rules.
//
// Rules:
// 1. Entries must not be empty
// 2. Each entry must have an account
// 3. Total debits must equal total credits
//
// TODO: Implement this function (BONUS - harder!)
func ValidateGLMap(entries []GLEntry) error {
	// YOUR CODE HERE
	// Hint:
	// 1. Check if entries is empty → return an error
	// 2. Loop through and check each has an Account → return error if empty
	// 3. Check IsBalanced() → return error if not balanced
	// 4. If all checks pass, return nil

	return nil // Replace this
}
