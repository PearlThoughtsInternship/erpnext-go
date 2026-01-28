package validation

import (
	"math"
	"testing"
)

func TestTotalDebit(t *testing.T) {
	tests := []struct {
		name     string
		entries  []GLEntry
		expected float64
	}{
		{
			name:     "single_entry",
			entries:  []GLEntry{{Debit: 100}},
			expected: 100,
		},
		{
			name: "multiple_entries",
			entries: []GLEntry{
				{Debit: 100},
				{Debit: 50},
				{Debit: 25},
			},
			expected: 175,
		},
		{
			name:     "empty_slice",
			entries:  []GLEntry{},
			expected: 0,
		},
		{
			name: "mixed_debit_credit",
			entries: []GLEntry{
				{Debit: 100, Credit: 0},
				{Debit: 0, Credit: 50},
				{Debit: 75, Credit: 0},
			},
			expected: 175, // Only count debits
		},
		{
			name: "decimal_values",
			entries: []GLEntry{
				{Debit: 100.50},
				{Debit: 200.75},
			},
			expected: 301.25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TotalDebit(tt.entries)
			if math.Abs(got-tt.expected) > epsilon {
				t.Errorf("TotalDebit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTotalCredit(t *testing.T) {
	tests := []struct {
		name     string
		entries  []GLEntry
		expected float64
	}{
		{
			name:     "single_entry",
			entries:  []GLEntry{{Credit: 100}},
			expected: 100,
		},
		{
			name: "multiple_entries",
			entries: []GLEntry{
				{Credit: 100},
				{Credit: 50},
				{Credit: 25},
			},
			expected: 175,
		},
		{
			name:     "empty_slice",
			entries:  []GLEntry{},
			expected: 0,
		},
		{
			name: "mixed_debit_credit",
			entries: []GLEntry{
				{Debit: 100, Credit: 0},
				{Debit: 0, Credit: 50},
				{Debit: 0, Credit: 75},
			},
			expected: 125, // Only count credits
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TotalCredit(tt.entries)
			if math.Abs(got-tt.expected) > epsilon {
				t.Errorf("TotalCredit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsBalanced(t *testing.T) {
	tests := []struct {
		name     string
		entries  []GLEntry
		expected bool
	}{
		{
			name:     "empty_slice_is_balanced",
			entries:  []GLEntry{},
			expected: true,
		},
		{
			name: "simple_balanced",
			entries: []GLEntry{
				{Account: "Debtors", Debit: 100},
				{Account: "Sales", Credit: 100},
			},
			expected: true,
		},
		{
			name: "multi_entry_balanced",
			entries: []GLEntry{
				{Account: "Debtors", Debit: 11800},
				{Account: "Sales", Credit: 10000},
				{Account: "CGST", Credit: 900},
				{Account: "SGST", Credit: 900},
			},
			expected: true,
		},
		{
			name: "unbalanced_more_debit",
			entries: []GLEntry{
				{Debit: 100},
				{Credit: 90},
			},
			expected: false,
		},
		{
			name: "unbalanced_more_credit",
			entries: []GLEntry{
				{Debit: 50},
				{Credit: 100},
			},
			expected: false,
		},
		{
			name: "within_epsilon_tolerance",
			entries: []GLEntry{
				{Debit: 100.00001},
				{Credit: 100.00002},
			},
			expected: true, // Difference is within epsilon
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBalanced(tt.entries)
			if got != tt.expected {
				t.Errorf("IsBalanced() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name     string
		entries  []GLEntry
		expected float64
	}{
		{
			name:     "empty_slice",
			entries:  []GLEntry{},
			expected: 0,
		},
		{
			name: "balanced_is_zero",
			entries: []GLEntry{
				{Debit: 100},
				{Credit: 100},
			},
			expected: 0,
		},
		{
			name: "more_debits",
			entries: []GLEntry{
				{Debit: 100},
				{Credit: 80},
			},
			expected: 20,
		},
		{
			name: "more_credits",
			entries: []GLEntry{
				{Debit: 50},
				{Credit: 100},
			},
			expected: -50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Difference(tt.entries)
			if math.Abs(got-tt.expected) > epsilon {
				t.Errorf("Difference() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestValidateGLMap tests the complete validation function (BONUS)
func TestValidateGLMap(t *testing.T) {
	tests := []struct {
		name      string
		entries   []GLEntry
		wantError bool
	}{
		{
			name:      "empty_entries_error",
			entries:   []GLEntry{},
			wantError: true,
		},
		{
			name: "missing_account_error",
			entries: []GLEntry{
				{Account: "Debtors", Debit: 100},
				{Account: "", Credit: 100}, // Missing account!
			},
			wantError: true,
		},
		{
			name: "unbalanced_error",
			entries: []GLEntry{
				{Account: "Debtors", Debit: 100},
				{Account: "Sales", Credit: 90},
			},
			wantError: true,
		},
		{
			name: "valid_entries",
			entries: []GLEntry{
				{Account: "Debtors", Debit: 100},
				{Account: "Sales", Credit: 100},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGLMap(tt.entries)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateGLMap() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestRealisticScenario tests with real-world ERPNext data
func TestRealisticScenario(t *testing.T) {
	// Sales Invoice with 18% GST (CGST 9% + SGST 9%)
	salesInvoice := []GLEntry{
		{Account: "Debtors - ACME", Debit: 11800},
		{Account: "Sales - ACME", Credit: 10000},
		{Account: "CGST Payable - ACME", Credit: 900},
		{Account: "SGST Payable - ACME", Credit: 900},
	}

	t.Run("sales_invoice_totals", func(t *testing.T) {
		if TotalDebit(salesInvoice) != 11800 {
			t.Errorf("Expected total debit 11800")
		}
		if TotalCredit(salesInvoice) != 11800 {
			t.Errorf("Expected total credit 11800")
		}
	})

	t.Run("sales_invoice_balanced", func(t *testing.T) {
		if !IsBalanced(salesInvoice) {
			t.Error("Sales invoice should be balanced")
		}
	})

	t.Run("sales_invoice_valid", func(t *testing.T) {
		if err := ValidateGLMap(salesInvoice); err != nil {
			t.Errorf("Sales invoice should be valid: %v", err)
		}
	})

	t.Logf("✅ Sales Invoice validated: Dr %.2f = Cr %.2f",
		TotalDebit(salesInvoice), TotalCredit(salesInvoice))
}
