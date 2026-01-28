package glentry

import (
	"testing"
	"time"
)

// TestGLEntryFields verifies that all required fields exist with correct types.
// If you've completed all TODOs correctly, this test will pass.
func TestGLEntryFields(t *testing.T) {
	// Create a GL entry with all fields populated
	entry := GLEntry{
		Name:         "GL-00001",
		Account:      "Sales - ACME",
		Company:      "ACME Corp",
		Debit:        0,
		Credit:       10000.00,
		VoucherType:  "Sales Invoice",
		VoucherNo:    "SINV-2024-00001",
		PartyType:    "Customer",
		Party:        "Acme Corporation",
		PostingDate:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		CreationDate: time.Now(),
		CostCenter:   "Main - ACME",
		IsCancelled:  false,
		IsOpening:    "No",
	}

	// Test that fields are accessible and have correct values
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Name", entry.Name, "GL-00001"},
		{"Account", entry.Account, "Sales - ACME"},
		{"Company", entry.Company, "ACME Corp"},
		{"Debit", entry.Debit, 0.0},
		{"Credit", entry.Credit, 10000.00},
		{"VoucherType", entry.VoucherType, "Sales Invoice"},
		{"VoucherNo", entry.VoucherNo, "SINV-2024-00001"},
		{"PartyType", entry.PartyType, "Customer"},
		{"Party", entry.Party, "Acme Corporation"},
		{"CostCenter", entry.CostCenter, "Main - ACME"},
		{"IsCancelled", entry.IsCancelled, false},
		{"IsOpening", entry.IsOpening, "No"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s: got %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}

	// Verify PostingDate is set correctly
	if entry.PostingDate.Year() != 2024 || entry.PostingDate.Month() != 1 {
		t.Errorf("PostingDate: got %v, want 2024-01-15", entry.PostingDate)
	}
}

// TestGLEntryDefaults verifies the NewGLEntry constructor sets proper defaults.
func TestGLEntryDefaults(t *testing.T) {
	entry := NewGLEntry("Debtors - ACME", 11800.00, 0)

	// Check that constructor sets values correctly
	if entry.Account != "Debtors - ACME" {
		t.Errorf("Account: got %q, want %q", entry.Account, "Debtors - ACME")
	}

	if entry.Debit != 11800.00 {
		t.Errorf("Debit: got %v, want %v", entry.Debit, 11800.00)
	}

	if entry.Credit != 0 {
		t.Errorf("Credit: got %v, want %v", entry.Credit, 0.0)
	}

	// Check defaults
	if entry.IsCancelled != false {
		t.Errorf("IsCancelled should default to false")
	}

	if entry.IsOpening != "No" {
		t.Errorf("IsOpening should default to \"No\", got %q", entry.IsOpening)
	}

	// PostingDate should be set (not zero)
	if entry.PostingDate.IsZero() {
		t.Error("PostingDate should be set to current time")
	}

	// CreationDate should be set (not zero)
	if entry.CreationDate.IsZero() {
		t.Error("CreationDate should be set to current time")
	}
}

// TestGLEntryIsValid tests the validation function.
func TestGLEntryIsValid(t *testing.T) {
	tests := []struct {
		name    string
		entry   GLEntry
		isValid bool
	}{
		{
			name:    "valid_debit_entry",
			entry:   GLEntry{Account: "Debtors - ACME", Debit: 100, Credit: 0},
			isValid: true,
		},
		{
			name:    "valid_credit_entry",
			entry:   GLEntry{Account: "Sales - ACME", Debit: 0, Credit: 100},
			isValid: true,
		},
		{
			name:    "invalid_both_debit_and_credit",
			entry:   GLEntry{Account: "Test", Debit: 100, Credit: 50},
			isValid: false,
		},
		{
			name:    "invalid_empty_account",
			entry:   GLEntry{Account: "", Debit: 100, Credit: 0},
			isValid: false,
		},
		{
			name:    "invalid_zero_amounts",
			entry:   GLEntry{Account: "Test", Debit: 0, Credit: 0},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entry.IsValid(); got != tt.isValid {
				t.Errorf("IsValid(): got %v, want %v", got, tt.isValid)
			}
		})
	}
}

// TestRealisticSalesInvoice creates entries like ERPNext would.
// This tests your understanding of how the fields work together.
func TestRealisticSalesInvoice(t *testing.T) {
	// A Sales Invoice with GST creates these entries:
	// Debtors (Dr) 11,800
	// Sales (Cr) 10,000
	// CGST (Cr) 900
	// SGST (Cr) 900

	entries := []GLEntry{
		{
			Account:     "Debtors - ACME",
			Debit:       11800.00,
			Credit:      0,
			VoucherType: "Sales Invoice",
			VoucherNo:   "SINV-2024-00001",
			PartyType:   "Customer",
			Party:       "Acme Corporation",
		},
		{
			Account:     "Sales - ACME",
			Debit:       0,
			Credit:      10000.00,
			VoucherType: "Sales Invoice",
			VoucherNo:   "SINV-2024-00001",
		},
		{
			Account:     "CGST Payable - ACME",
			Debit:       0,
			Credit:      900.00,
			VoucherType: "Sales Invoice",
			VoucherNo:   "SINV-2024-00001",
		},
		{
			Account:     "SGST Payable - ACME",
			Debit:       0,
			Credit:      900.00,
			VoucherType: "Sales Invoice",
			VoucherNo:   "SINV-2024-00001",
		},
	}

	// Calculate totals
	var totalDebit, totalCredit float64
	for _, e := range entries {
		totalDebit += e.Debit
		totalCredit += e.Credit
	}

	// The golden rule: Debits must equal Credits
	if totalDebit != totalCredit {
		t.Errorf("Entries don't balance! Debit: %v, Credit: %v", totalDebit, totalCredit)
	}

	// Verify the totals
	if totalDebit != 11800.00 {
		t.Errorf("Total debit: got %v, want %v", totalDebit, 11800.00)
	}

	if totalCredit != 11800.00 {
		t.Errorf("Total credit: got %v, want %v", totalCredit, 11800.00)
	}

	t.Logf("✅ Sales Invoice balanced: Dr %.2f = Cr %.2f", totalDebit, totalCredit)
}
